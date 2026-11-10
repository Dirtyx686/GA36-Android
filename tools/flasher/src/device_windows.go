//go:build windows

package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

var shell = syscall.NewLazyDLL("shell32.dll")

func ensureAdmin() error {
	v, _, _ := shell.NewProc("IsUserAnAdmin").Call()
	if v != 0 {
		return nil
	}
	exe, e := os.Executable()
	if e != nil {
		return e
	}
	verb, _ := syscall.UTF16PtrFromString("runas")
	file, _ := syscall.UTF16PtrFromString(exe)
	args, _ := syscall.UTF16PtrFromString("--elevated")
	dir, _ := syscall.UTF16PtrFromString(filepath.Dir(exe))
	r, _, _ := shell.NewProc("ShellExecuteW").Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(file)), uintptr(unsafe.Pointer(args)), uintptr(unsafe.Pointer(dir)), 1)
	if r <= 32 {
		return errors.New("Administrator-Anfrage abgebrochen")
	}
	os.Exit(0)
	return nil
}
func ps(script string) ([]byte, error) {
	p := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", "$ErrorActionPreference='Stop'; [Console]::OutputEncoding=[System.Text.Encoding]::UTF8; "+script)
	p.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	b, e := p.Output()
	if e != nil {
		return nil, fmt.Errorf("Windows-Datenträgerabfrage fehlgeschlagen: %w", e)
	}
	return []byte(strings.TrimPrefix(strings.TrimSpace(string(b)), "\ufeff")), nil
}

type winDisk struct {
	Number                                                              int
	FriendlyName, SerialNumber, UniqueId, Path, BusType, PartitionStyle string
	Size                                                                int64
	IsBoot, IsSystem, IsReadOnly                                        bool
}

func listDisks(paths ...string) ([]Disk, error) {
	b, e := ps("ConvertTo-Json -Depth 4 -Compress -InputObject @(Get-Disk | Select-Object Number,FriendlyName,SerialNumber,UniqueId,Path,@{N='BusType';E={$_.BusType.ToString()}},@{N='PartitionStyle';E={$_.PartitionStyle.ToString()}},Size,IsBoot,IsSystem,IsReadOnly)")
	if e != nil {
		return nil, e
	}
	var in []winDisk
	if e = json.Unmarshal(b, &in); e != nil {
		return nil, e
	}
	var out []Disk
	for _, x := range in {
		if x.IsBoot || x.IsSystem || x.IsReadOnly || (x.BusType != "USB" && x.BusType != "SD" && x.BusType != "MMC") {
			continue
		}
		safe := true
		// Image and executable drive letters are resolved to physical disk numbers.
		for _, p := range paths {
			vol := filepath.VolumeName(p)
			if len(vol) != 2 || vol[1] != ':' {
				return nil, errors.New("Paket bitte auf ein lokales Laufwerk entpacken")
			}
			letter := strings.ToUpper(vol[:1])
			if letter[0] < 'A' || letter[0] > 'Z' {
				return nil, errors.New("Unbekanntes Paketlaufwerk")
			}
			q, e := ps("ConvertTo-Json -Compress -InputObject @((Get-Partition -DriveLetter '" + letter + "').DiskNumber)")
			if e != nil {
				return nil, e
			}
			var nums []int
			if e = json.Unmarshal(q, &nums); e != nil {
				return nil, e
			}
			for _, n := range nums {
				if n == x.Number {
					safe = false
				}
			}
		}
		id := fmt.Sprintf("%d|%s|%s|%s|%d|%s", x.Number, x.Path, x.UniqueId, x.SerialNumber, x.Size, x.PartitionStyle)
		out = append(out, Disk{Path: `\\.\PhysicalDrive` + strconv.Itoa(x.Number), ID: strconv.Itoa(x.Number), Model: x.FriendlyName, Serial: x.SerialNumber, Bytes: x.Size, Bus: x.BusType, Identity: id, Safe: safe, Raw: x.PartitionStyle == "RAW"})
	}
	return out, nil
}
func winOpen(path string, access, flags uint32) (*os.File, error) {
	p, e := syscall.UTF16PtrFromString(path)
	if e != nil {
		return nil, e
	}
	h, e := syscall.CreateFile(p, access, syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE, nil, syscall.OPEN_EXISTING, flags, 0)
	if e != nil {
		return nil, e
	}
	return os.NewFile(uintptr(h), path), nil
}
func winIo(f *os.File, code uint32, out []byte) error {
	var p *byte
	if len(out) > 0 {
		p = &out[0]
	}
	var n uint32
	return syscall.DeviceIoControl(syscall.Handle(f.Fd()), code, nil, 0, p, uint32(len(out)), &n, nil)
}
func winSize(f *os.File) (int64, error) {
	var b [8]byte
	e := winIo(f, 0x7405c, b[:])
	return int64(binary.LittleEndian.Uint64(b[:])), e
}
func openDevice(d Disk, paths ...string) (*Device, error) {
	if !d.Safe {
		return nil, errors.New("Geschütztes Laufwerk")
	}
	n, e := strconv.Atoi(d.ID)
	if e != nil || n < 0 || n > 999 {
		return nil, errors.New("Ungültige Datenträgernummer")
	}
	var volumes []string
	if !d.Raw {
		raw, e := ps(`ConvertTo-Json -Compress -InputObject @((Get-Partition -DiskNumber ` + strconv.Itoa(n) + `).AccessPaths | Where-Object { $_ -and $_.StartsWith('\\?\Volume{') })`)
		if e != nil {
			return nil, e
		}
		if e = json.Unmarshal(raw, &volumes); e != nil {
			return nil, e
		}
	}
	var locks []*os.File
	cleanup := func() {
		for _, f := range locks {
			winIo(f, 0x9001c, nil)
			f.Close()
		}
	}
	for _, v := range volumes {
		if !strings.HasPrefix(v, `\\?\Volume{`) || !strings.HasSuffix(v, `}\`) {
			cleanup()
			return nil, errors.New("Unbekannter Volumepfad")
		}
		f, e := winOpen(strings.TrimSuffix(v, `\`), syscall.GENERIC_READ|syscall.GENERIC_WRITE, 0)
		if e != nil {
			cleanup()
			return nil, e
		}
		locks = append(locks, f)
		if e = winIo(f, 0x90018, nil); e != nil {
			cleanup()
			return nil, fmt.Errorf("Karte ist noch in Benutzung: %w", e)
		}
		if e = winIo(f, 0x90020, nil); e != nil {
			cleanup()
			return nil, e
		}
	}
	fresh, e := listDisks(paths...)
	if e != nil {
		cleanup()
		return nil, e
	}
	found := false
	for _, x := range fresh {
		if x.Path == d.Path && x.Identity == d.Identity && x.Safe {
			found = true
		}
	}
	if !found {
		cleanup()
		return nil, errors.New("Ziel hat sich geändert")
	}
	f, e := winOpen(d.Path, syscall.GENERIC_READ|syscall.GENERIC_WRITE, 0x80000000|0x20000000)
	if e != nil {
		cleanup()
		return nil, e
	}
	s, e := winSize(f)
	if e != nil || s != d.Bytes {
		f.Close()
		cleanup()
		return nil, errors.New("Kartengröße hat sich geändert")
	}
	var once sync.Once
	var closed error
	closeFn := func() error { once.Do(func() { closed = f.Close(); cleanup() }); return closed }
	guard := func() error {
		size, e := winSize(f)
		if e != nil || size != d.Bytes {
			return errors.New("Karte getrennt oder geändert")
		}
		return nil
	}
	return &Device{File: f, Guard: guard, Close: closeFn, Finish: func() error { return f.Sync() }, ReopenRead: func() (*os.File, error) {
		if e := f.Sync(); e != nil {
			return nil, e
		}
		return winOpen(d.Path, syscall.GENERIC_READ, 0x20000000)
	}}, nil
}
func pauseExit() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "--check") {
		return
	}
	fmt.Print("\nEnter zum Schließen …")
	var s string
	fmt.Scanln(&s)
}

func listInboxes(exe string) ([]Inbox, error) {
	disks, e := listDisks(exe)
	if e != nil {
		return nil, e
	}
	safe := map[int]Disk{}
	for _, d := range disks {
		n, _ := strconv.Atoi(d.ID)
		if d.Safe {
			safe[n] = d
		}
	}
	b, e := ps(`$r=@(); foreach($v in @(Get-Volume | Where-Object {$_.FileSystemLabel -eq 'GA36_APPS' -and $_.FileSystem -eq 'FAT32'})){if($v.DriveLetter){$p=Get-Partition -DriveLetter $v.DriveLetter; $r+= [pscustomobject]@{Root=($v.DriveLetter.ToString()+':\');Disk=[int]$p.DiskNumber;Offset=[long]$p.Offset;Size=[long]$p.Size}}};ConvertTo-Json -Compress -InputObject @($r)`)
	if e != nil {
		return nil, e
	}
	var rows []struct {
		Root         string
		Disk         int
		Offset, Size int64
	}
	if e = json.Unmarshal(b, &rows); e != nil {
		return nil, e
	}
	var out []Inbox
	for _, r := range rows {
		d, ok := safe[r.Disk]
		if !ok || r.Offset != 1732268032 || r.Size != 4294967296 {
			continue
		}
		root := r.Root
		out = append(out, Inbox{Description: d.Model + " · " + d.Path + " · " + root, Open: func() (string, func(), error) { return root, func() {}, nil }})
	}
	return out, nil
}
func inboxFree(root string) (uint64, error) {
	p, e := syscall.UTF16PtrFromString(root)
	if e != nil {
		return 0, e
	}
	var n uint64
	r, _, err := syscall.NewLazyDLL("kernel32.dll").NewProc("GetDiskFreeSpaceExW").Call(uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(&n)), 0, 0)
	if r == 0 {
		return 0, err
	}
	return n, nil
}
func syncInbox(root string) error {
	volume := filepath.VolumeName(root)
	if len(volume) != 2 {
		return errors.New("Unbekannter Volumepfad")
	}
	f, e := winOpen(`\\.\`+volume, syscall.GENERIC_READ|syscall.GENERIC_WRITE, 0)
	if e != nil {
		return e
	}
	defer f.Close()
	return f.Sync()
}

func init() {
	k := syscall.NewLazyDLL("kernel32.dll")
	k.NewProc("SetConsoleOutputCP").Call(65001)
	k.NewProc("SetConsoleCP").Call(65001)
}
