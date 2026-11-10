//go:build linux

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

type lsNode struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Size     int64     `json:"size"`
	Model    string    `json:"model"`
	Serial   string    `json:"serial"`
	Tran     string    `json:"tran"`
	RM       bool      `json:"rm"`
	RO       bool      `json:"ro"`
	Mounts   []*string `json:"mountpoints"`
	Children []lsNode  `json:"children"`
}

func ensureAdmin() error {
	if os.Geteuid() == 0 {
		return nil
	}
	exe, e := os.Executable()
	if e != nil {
		return e
	}
	sudo, e := exec.LookPath("sudo")
	if e != nil {
		return errors.New("Bitte im Terminal mit sudo starten")
	}
	return syscall.Exec(sudo, []string{"sudo", "--", exe, "--console"}, os.Environ())
}
func below(path, root string) bool {
	root = filepath.Clean(root)
	return path == root || strings.HasPrefix(path, root+string(os.PathSeparator)) || root == "/"
}
func protectedMount(m string, paths ...string) bool {
	for _, x := range []string{"/", "/boot", "/boot/efi", "/home", "/usr", "/var", "/opt"} {
		if m == x {
			return true
		}
	}
	for _, p := range paths {
		if below(p, m) {
			return true
		}
	}
	return false
}
func flatten(n lsNode) []string {
	var out []string
	for _, p := range n.Mounts {
		if p != nil && *p != "" {
			out = append(out, *p)
		}
	}
	for _, c := range n.Children {
		out = append(out, flatten(c)...)
	}
	return out
}
func listDisks(paths ...string) ([]Disk, error) {
	out, e := exec.Command("lsblk", "--json", "--bytes", "--paths", "--output", "NAME,TYPE,SIZE,MODEL,SERIAL,TRAN,RM,RO,MOUNTPOINTS").Output()
	if e != nil {
		return nil, fmt.Errorf("Kartenliste konnte nicht gelesen werden: %w", e)
	}
	var tree struct {
		Nodes []lsNode `json:"blockdevices"`
	}
	if e = json.Unmarshal(out, &tree); e != nil {
		return nil, e
	}
	var result []Disk
	for _, n := range tree.Nodes {
		if n.Type != "disk" || n.RO || !strings.HasPrefix(n.Name, "/dev/") {
			continue
		}
		base := filepath.Base(n.Name)
		sys := filepath.Join("/sys/class/block", base)
		t, _ := os.ReadFile(filepath.Join(sys, "device/type"))
		isSD := strings.TrimSpace(string(t)) == "SD"
		if n.Tran != "usb" && !isSD {
			continue
		}
		seq, e := os.ReadFile(filepath.Join(sys, "diskseq"))
		if e != nil {
			continue
		}
		mounts := flatten(n)
		safe := true
		for _, m := range mounts {
			if protectedMount(m, paths...) {
				safe = false
			}
		}
		id := fmt.Sprintf("%s|%s|%s|%d|%s", n.Name, n.Model, n.Serial, n.Size, strings.TrimSpace(string(seq)))
		result = append(result, Disk{Path: n.Name, Model: strings.TrimSpace(n.Model), Serial: strings.TrimSpace(n.Serial), Bytes: n.Size, Bus: n.Tran, Seq: strings.TrimSpace(string(seq)), Identity: id, Safe: safe, Mounts: mounts})
	}
	return result, nil
}
func ioctl(fd uintptr, code uintptr, arg unsafe.Pointer) error {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, code, uintptr(arg))
	if e != 0 {
		return e
	}
	return nil
}
func openDevice(d Disk, paths ...string) (*Device, error) {
	if !d.Safe {
		return nil, errors.New("Geschütztes Laufwerk")
	}
	sort.Slice(d.Mounts, func(i, j int) bool { return len(d.Mounts[i]) > len(d.Mounts[j]) })
	seen := map[string]bool{}
	for _, m := range d.Mounts {
		if seen[m] {
			continue
		}
		seen[m] = true
		if protectedMount(m, paths...) {
			return nil, errors.New("Image oder System liegt auf dem Ziel")
		}
		if out, e := exec.Command("umount", "--", m).CombinedOutput(); e != nil {
			return nil, fmt.Errorf("Karte ist noch in Benutzung: %s (%w)", out, e)
		}
	}
	fresh, e := listDisks(paths...)
	if e != nil {
		return nil, e
	}
	found := false
	for _, x := range fresh {
		if x.Path == d.Path && x.Identity == d.Identity && x.Safe && len(x.Mounts) == 0 {
			found = true
		}
	}
	if !found {
		return nil, errors.New("Karte geändert oder noch eingehängt")
	}
	f, e := os.OpenFile(d.Path, os.O_RDWR|syscall.O_EXCL|syscall.O_NOFOLLOW|syscall.O_SYNC, 0)
	if e != nil {
		return nil, e
	}
	ok := false
	defer func() {
		if !ok {
			f.Close()
		}
	}()
	st, e := f.Stat()
	if e != nil {
		return nil, e
	}
	if st.Mode()&os.ModeDevice == 0 || st.Mode()&os.ModeCharDevice != 0 {
		return nil, errors.New("Ziel ist kein Blockgerät")
	}
	var size uint64
	if e = ioctl(f.Fd(), 0x80081272, unsafe.Pointer(&size)); e != nil || size != uint64(d.Bytes) {
		return nil, errors.New("Kartengröße hat sich geändert")
	}
	guard := func() error {
		b, e := os.ReadFile(filepath.Join("/sys/class/block", filepath.Base(d.Path), "diskseq"))
		if e != nil || strings.TrimSpace(string(b)) != d.Seq {
			return errors.New("Karte wurde getrennt oder ausgetauscht")
		}
		var s uint64
		if e = ioctl(f.Fd(), 0x80081272, unsafe.Pointer(&s)); e != nil || s != size {
			return errors.New("Datenträger nicht mehr unverändert erreichbar")
		}
		return nil
	}
	var once sync.Once
	var closeErr error
	closeFn := func() error { once.Do(func() { closeErr = f.Close() }); return closeErr }
	dev := &Device{File: f, Guard: guard, Close: closeFn, Finish: func() error {
		if e := f.Sync(); e != nil {
			return e
		}
		return ioctl(f.Fd(), 0x1261, nil)
	}, ReopenRead: func() (*os.File, error) {
		if e := f.Sync(); e != nil {
			return nil, e
		}
		if e := ioctl(f.Fd(), 0x1261, nil); e != nil {
			return nil, e
		}
		return os.OpenFile(d.Path, os.O_RDONLY|syscall.O_DIRECT|syscall.O_NOFOLLOW, 0)
	}}
	ok = true
	return dev, nil
}

var _ = strconv.IntSize

func pauseExit() {}

type fsNode struct {
	Name, Type, Fstype, Label, UUID string
	Start                           int64
	Size                            int64
	Mounts                          []*string `json:"mountpoints"`
	Children                        []fsNode  `json:"children"`
}

func listInboxes(exe string) ([]Inbox, error) {
	disks, e := listDisks(exe)
	if e != nil {
		return nil, e
	}
	safe := map[string]Disk{}
	for _, d := range disks {
		if d.Safe {
			safe[d.Path] = d
		}
	}
	b, e := exec.Command("lsblk", "--json", "--bytes", "--paths", "--output", "NAME,TYPE,FSTYPE,LABEL,UUID,START,SIZE,MOUNTPOINTS").Output()
	if e != nil {
		return nil, e
	}
	var tree struct {
		Nodes []fsNode `json:"blockdevices"`
	}
	if e = json.Unmarshal(b, &tree); e != nil {
		return nil, e
	}
	var out []Inbox
	var walk func(fsNode, Disk)
	walk = func(n fsNode, d Disk) {
		if n.Type == "part" && n.Start == 3383336 && n.Size == 4294967296 && strings.EqualFold(n.UUID, "4736-0701") && n.Label == "GA36_APPS" && n.Fstype == "vfat" {
			node := n
			disk := d
			out = append(out, Inbox{Description: disk.Model + " · " + disk.Path + " · GA36_APPS", Open: func() (string, func(), error) {
				seq, e := os.ReadFile(filepath.Join("/sys/class/block", filepath.Base(disk.Path), "diskseq"))
				if e != nil || strings.TrimSpace(string(seq)) != disk.Seq {
					return "", nil, errors.New("Karte wurde gewechselt")
				}
				for _, m := range node.Mounts {
					if m != nil && *m != "" {
						return *m, func() {}, nil
					}
				}
				root, e := os.MkdirTemp("", "ga36-inbox-")
				if e != nil {
					return "", nil, e
				}
				b, e := exec.Command("mount", "-t", "vfat", "-o", "rw,nosuid,nodev,noexec,umask=077", "--", node.Name, root).CombinedOutput()
				if e != nil {
					os.Remove(root)
					return "", nil, fmt.Errorf("Austauschbereich kann nicht geöffnet werden: %s %w", b, e)
				}
				cleanup := func() {
					if b, e := exec.Command("umount", "--", root).CombinedOutput(); e != nil {
						fmt.Fprintln(os.Stderr, "Bitte Karte vor dem Entnehmen über das System auswerfen:", string(b))
					} else {
						os.Remove(root)
					}
				}
				return root, cleanup, nil
			}})
		}
		for _, c := range n.Children {
			walk(c, d)
		}
	}
	for _, n := range tree.Nodes {
		if d, ok := safe[n.Name]; ok {
			walk(n, d)
		}
	}
	return out, nil
}
func inboxFree(root string) (uint64, error) {
	var s syscall.Statfs_t
	e := syscall.Statfs(root, &s)
	return uint64(s.Bavail) * uint64(s.Bsize), e
}
func syncInbox(root string) error {
	f, e := os.Open(root)
	if e != nil {
		return e
	}
	defer f.Close()
	e = f.Sync()
	if errors.Is(e, syscall.EINVAL) {
		return nil
	}
	return e
}
