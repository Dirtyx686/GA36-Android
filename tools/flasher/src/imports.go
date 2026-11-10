package main

import (
	"archive/zip"
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
)

const inboxMarker = "GA36_ANDROID7_R01_APK_INBOX\n"

type Inbox struct {
	Root, Description string
	Open              func() (string, func(), error)
}
type APKInfo struct {
	Package, Version string
	MinSDK           int
	Native           []string
	Bytes            int64
	SHA              string
}

var packageRE = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*(\.[A-Za-z][A-Za-z0-9_]*)+$`)
var obbRE = regexp.MustCompile(`^(main|patch)\.([0-9]+)\.([A-Za-z][A-Za-z0-9_]*(?:\.[A-Za-z][A-Za-z0-9_]*)+)\.obb$`)

func le16(b []byte) uint16 { return binary.LittleEndian.Uint16(b) }
func le32(b []byte) uint32 { return binary.LittleEndian.Uint32(b) }
func stringPool(b []byte) ([]string, error) {
	if len(b) < 28 {
		return nil, errors.New("short stringpool")
	}
	count := int(le32(b[8:]))
	start := int(le32(b[20:]))
	hdr := int(le16(b[2:]))
	flags := le32(b[16:])
	if count > 100000 || hdr < 28 || hdr+count*4 > len(b) || start > len(b) {
		return nil, errors.New("bad stringpool")
	}
	out := make([]string, count)
	for i := range out {
		at := start + int(le32(b[hdr+i*4:]))
		if at < start || at >= len(b) {
			return nil, errors.New("bad string offset")
		}
		if flags&0x100 != 0 {
			readLen := func() (int, error) {
				if at >= len(b) {
					return 0, io.ErrUnexpectedEOF
				}
				v := int(b[at])
				at++
				if v&0x80 != 0 {
					if at >= len(b) {
						return 0, io.ErrUnexpectedEOF
					}
					v = (v&0x7f)<<8 | int(b[at])
					at++
				}
				return v, nil
			}
			if _, e := readLen(); e != nil {
				return nil, e
			}
			n, e := readLen()
			if e != nil || n > len(b)-at-1 {
				return nil, errors.New("bad UTF8 length")
			}
			out[i] = string(b[at : at+n])
		} else {
			if at+2 > len(b) {
				return nil, io.ErrUnexpectedEOF
			}
			n := uint32(le16(b[at:]))
			at += 2
			if n&0x8000 != 0 {
				if at+2 > len(b) {
					return nil, io.ErrUnexpectedEOF
				}
				n = (n&0x7fff)<<16 | uint32(le16(b[at:]))
				at += 2
			}
			if uint64(n)*2 > uint64(len(b)-at) {
				return nil, errors.New("bad UTF16 length")
			}
			v := make([]uint16, n)
			for j := range v {
				v[j] = le16(b[at+j*2:])
			}
			out[i] = string(utf16.Decode(v))
		}
	}
	return out, nil
}
func parseManifest(b []byte) (APKInfo, error) {
	info := APKInfo{MinSDK: 1}
	bad := errors.New("AndroidManifest.xml nicht lesbar")
	if len(b) < 8 || len(b) > 4*1024*1024 || le16(b) != 3 || int(le32(b[4:])) != len(b) {
		return info, bad
	}
	var pool []string
	got := false
	get := func(n uint32) (string, error) {
		if n >= uint32(len(pool)) {
			return "", bad
		}
		return pool[n], nil
	}
	for off := int(le16(b[2:])); off < len(b); {
		if off < 8 || off+8 > len(b) {
			return info, bad
		}
		size := int(le32(b[off+4:]))
		if size < 8 || size > len(b)-off {
			return info, bad
		}
		c := b[off : off+size]
		switch le16(c) {
		case 1:
			var e error
			pool, e = stringPool(c)
			if e != nil {
				return info, e
			}
		case 0x102:
			if len(c) < 36 {
				return info, bad
			}
			tag, e := get(le32(c[20:]))
			if e != nil {
				return info, e
			}
			start := 16 + int(le16(c[24:]))
			stride := int(le16(c[26:]))
			count := int(le16(c[28:]))
			if stride < 20 || start < 36 || count > (len(c)-start)/stride || start > len(c) {
				return info, bad
			}
			for j := 0; j < count; j++ {
				a := c[start+j*stride:]
				name, e := get(le32(a[4:]))
				if e != nil {
					return info, e
				}
				raw := le32(a[8:])
				typ := a[15]
				value := ""
				number := int(le32(a[16:]))
				if raw != 0xffffffff {
					value, e = get(raw)
					if e != nil {
						return info, e
					}
				} else if typ == 3 {
					value, e = get(le32(a[16:]))
					if e != nil {
						return info, e
					}
				} else if typ == 0x10 || typ == 0x11 {
					value = strconv.Itoa(number)
				}
				if tag == "manifest" {
					switch name {
					case "package":
						info.Package = value
						got = true
					case "versionName":
						info.Version = value
					case "split":
						if value != "" {
							return info, errors.New("Split-APK nicht unterstützt; vollständige Einzel-APK verwenden")
						}
					}
				}
				if tag == "uses-sdk" && name == "minSdkVersion" {
					v, e := strconv.Atoi(value)
					if e != nil || v < 1 {
						return info, errors.New("Unbekannte Mindest-Android-Version")
					}
					info.MinSDK = v
				}
			}
		}
		off += size
	}
	if !got || !packageRE.MatchString(info.Package) {
		return info, bad
	}
	return info, nil
}
func inspectAPK(path string) (APKInfo, error) {
	var info APKInfo
	f, e := os.Open(path)
	if e != nil {
		return info, e
	}
	defer f.Close()
	st, e := f.Stat()
	if e != nil || !st.Mode().IsRegular() || st.Size() <= 0 || st.Size() > 0xffffffff {
		return info, errors.New("APK muss eine reguläre Datei kleiner4GB sein")
	}
	z, e := zip.NewReader(f, st.Size())
	if e != nil {
		return info, e
	}
	found := false
	hasARM := false
	abis := map[string]bool{}
	for _, entry := range z.File {
		if entry.Name == "AndroidManifest.xml" {
			if found || entry.UncompressedSize64 > 4*1024*1024 {
				return info, errors.New("Ungültiges Manifest")
			}
			r, e := entry.Open()
			if e != nil {
				return info, e
			}
			b, e := io.ReadAll(io.LimitReader(r, 4*1024*1024+1))
			r.Close()
			if e != nil {
				return info, e
			}
			info, e = parseManifest(b)
			if e != nil {
				return info, e
			}
			found = true
		}
		if strings.HasPrefix(entry.Name, "lib/") {
			parts := strings.Split(entry.Name, "/")
			if len(parts) >= 3 && !entry.FileInfo().IsDir() {
				abis[parts[1]] = true
				if parts[1] == "armeabi-v7a" || parts[1] == "armeabi" {
					hasARM = true
				}
			}
		}
	}
	if !found {
		return info, errors.New("Kein Android-Manifest")
	}
	if info.MinSDK > 24 {
		return info, fmt.Errorf("APK braucht Android-API%d; dieses System hat API24 (Android7.0)", info.MinSDK)
	}
	if len(abis) > 0 && !hasARM {
		return info, errors.New("APK enthält keine ARM-32-Bit-Version")
	}
	for a := range abis {
		info.Native = append(info.Native, a)
	}
	info.Bytes = st.Size()
	info.SHA, e = hashFile(f)
	return info, e
}
func jobID() (string, error) {
	b := make([]byte, 16)
	_, e := rand.Read(b)
	return hex.EncodeToString(b), e
}
func checkInbox(root string) error {
	f, e := os.Open(filepath.Join(root, "GA36-APPS.TXT"))
	if e != nil {
		return e
	}
	b, e := io.ReadAll(io.LimitReader(f, int64(len(inboxMarker)+1)))
	f.Close()
	if e != nil || string(b) != inboxMarker {
		return errors.New("Keine gültige GA36-Austauschpartition")
	}
	for _, n := range []string{"APKS", "OBB"} {
		p := filepath.Join(root, n)
		s, e := os.Lstat(p)
		if e != nil || !s.IsDir() || s.Mode()&os.ModeSymlink != 0 {
			return errors.New("Ungültiger Austauschordner")
		}
	}
	return nil
}
func verifiedCopy(source, dest, wanted string) error {
	f, e := os.Open(source)
	if e != nil {
		return e
	}
	defer f.Close()
	st, e := f.Stat()
	if e != nil || !st.Mode().IsRegular() || st.Size() > 0xffffffff {
		return errors.New("Datei ist für FAT32 ungeeignet")
	}
	o, e := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		return e
	}
	ok := false
	defer func() {
		o.Close()
		if !ok {
			os.Remove(dest)
		}
	}()
	n, e := io.Copy(o, f)
	if e != nil || n != st.Size() {
		return errors.New("Kopieren fehlgeschlagen; Platz auf der Austauschpartition prüfen")
	}
	if e = o.Sync(); e != nil {
		return e
	}
	if e = o.Close(); e != nil {
		return e
	}
	r, e := os.Open(dest)
	if e != nil {
		return e
	}
	h, e := hashFile(r)
	r.Close()
	if e != nil || h != wanted {
		return errors.New("Prüfsumme nach dem Kopieren falsch")
	}
	ok = true
	return nil
}
func writeSmall(path, text string) error {
	f, e := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		return e
	}
	_, e = f.WriteString(text)
	if e == nil {
		e = f.Sync()
	}
	c := f.Close()
	if e != nil {
		return e
	}
	return c
}
func queueAPK(root, path string) (APKInfo, error) {
	info, e := inspectAPK(path)
	if e != nil {
		return info, e
	}
	if e = checkInbox(root); e != nil {
		return info, e
	}
	free, e := inboxFree(root)
	if e != nil || free < uint64(info.Bytes)+1048576 {
		return info, errors.New("Zu wenig Platz im Austauschbereich; Konsole starten oder alte Aufträge entfernen")
	}
	id, e := jobID()
	if e != nil {
		return info, e
	}
	final := filepath.Join(root, "APKS", info.SHA+"-"+id+".apk")
	tmp := final + ".part"
	if e = verifiedCopy(path, tmp, info.SHA); e != nil {
		return info, e
	}
	defer os.Remove(tmp)
	if e = writeSmall(final+".pkg", info.Package+"\n"); e != nil {
		return info, e
	}
	if e = os.Rename(tmp, final); e != nil {
		os.Remove(final + ".pkg")
		return info, e
	}
	return info, syncInbox(root)
}

type OBBFile struct {
	Path, Name, SHA string
	Bytes           int64
}

func inspectOBB(path string) (string, []OBBFile, error) {
	st, e := os.Lstat(path)
	if e != nil {
		return "", nil, e
	}
	var paths []string
	if st.IsDir() {
		entries, e := os.ReadDir(path)
		if e != nil {
			return "", nil, e
		}
		for _, x := range entries {
			if strings.HasSuffix(strings.ToLower(x.Name()), ".obb") {
				paths = append(paths, filepath.Join(path, x.Name()))
			}
		}
	} else {
		paths = []string{path}
	}
	if len(paths) == 0 {
		return "", nil, errors.New("Keine OBB-Dateien gefunden")
	}
	var pkg string
	var files []OBBFile
	for _, p := range paths {
		name := filepath.Base(p)
		parts := obbRE.FindStringSubmatch(name)
		if parts == nil {
			return "", nil, fmt.Errorf("OBB-Dateiname nicht standardgemäß: %s", name)
		}
		if pkg != "" && pkg != parts[3] {
			return "", nil, errors.New("Ordner enthält Daten verschiedener Apps")
		}
		pkg = parts[3]
		f, e := os.Open(p)
		if e != nil {
			return "", nil, e
		}
		s, e := f.Stat()
		if e != nil || !s.Mode().IsRegular() || s.Size() == 0 || s.Size() > 0xffffffff {
			f.Close()
			return "", nil, errors.New("OBB muss zwischen1Byte und4GB liegen")
		}
		h, e := hashFile(f)
		f.Close()
		if e != nil {
			return "", nil, e
		}
		files = append(files, OBBFile{p, name, h, s.Size()})
	}
	return pkg, files, nil
}
func queueOBB(root, path string) (string, error) {
	pkg, files, e := inspectOBB(path)
	if e != nil {
		return "", e
	}
	if e = checkInbox(root); e != nil {
		return "", e
	}
	var needed uint64
	for _, f := range files {
		needed += uint64(f.Bytes)
	}
	free, e := inboxFree(root)
	if e != nil || free < needed+1048576 {
		return "", errors.New("OBB-Daten passen nicht in den4GB-Austauschbereich; in getrennten Aufträgen übertragen")
	}
	id, e := jobID()
	if e != nil {
		return "", e
	}
	tmp := filepath.Join(root, "OBB", ".upload-"+id)
	if e = os.Mkdir(tmp, 0700); e != nil {
		return "", e
	}
	defer os.RemoveAll(tmp)
	manifest := ""
	for _, f := range files {
		if e = verifiedCopy(f.Path, filepath.Join(tmp, f.Name), f.SHA); e != nil {
			return "", e
		}
		manifest += f.SHA + " " + f.Name + "\n"
	}
	if e = writeSmall(filepath.Join(tmp, "PACKAGE.TXT"), pkg+"\n"); e != nil {
		return "", e
	}
	if e = writeSmall(filepath.Join(tmp, "SHA256.TXT"), manifest); e != nil {
		return "", e
	}
	if e = os.Rename(tmp, filepath.Join(root, "OBB", id)); e != nil {
		return "", e
	}
	return pkg, syncInbox(root)
}
func importMenu(exe string, in *bufio.Reader) error {
	if e := ensureAdmin(); e != nil {
		return e
	}
	volumes, e := listInboxes(exe)
	if e != nil {
		return e
	}
	if len(volumes) == 0 {
		return errors.New("Keine GA36_APPS-Partition gefunden. Zuerst das neue R01-Image auf eine Karte schreiben.")
	}
	for i, v := range volumes {
		fmt.Printf(" %d. %s\n", i+1, v.Description)
	}
	s, e := ask(in, "Zielkarte auswählen (Enter = Abbruch): ")
	if e != nil {
		return e
	}
	i, e := strconv.Atoi(s)
	if e != nil || i < 1 || i > len(volumes) {
		return errors.New("Abgebrochen")
	}
	root, close, e := volumes[i-1].Open()
	if e != nil {
		return e
	}
	defer close()
	if e = checkInbox(root); e != nil {
		return e
	}
	for {
		fmt.Println("\n1 APK vorbereiten\n2 OBB-Datei oder OBB-Ordner vorbereiten\n3 Letztes Import-Ergebnis anzeigen\n0 Zurück / Ende")
		s, e = ask(in, "Auswahl: ")
		if e != nil {
			return e
		}
		switch s {
		case "0", "":
			return nil
		case "1", "2":
			p, e := ask(in, "Datei- oder Ordnerpfad (auch hineinziehen möglich): ")
			if e != nil {
				return e
			}
			p = strings.Trim(strings.TrimSpace(p), "\"'")
			fmt.Println("Dateien und Kompatibilität werden geprüft …")
			if s == "1" {
				info, e := queueAPK(root, p)
				if e != nil {
					fmt.Println("Nicht übernommen:", e)
					continue
				}
				fmt.Println("APK bereit:", info.Package, "– Installation beim nächsten Start der Konsole.")
			} else {
				pkg, e := queueOBB(root, p)
				if e != nil {
					fmt.Println("Nicht übernommen:", e)
					continue
				}
				fmt.Println("OBB bereit:", pkg, "– Übernahme beim nächsten Start.")
			}
		case "3":
			f, e := os.Open(filepath.Join(root, "INSTALL.TXT"))
			if e != nil {
				fmt.Println(e)
				continue
			}
			st, e := f.Stat()
			if e != nil {
				f.Close()
				fmt.Println(e)
				continue
			}
			if st.Size() > 256*1024 {
				f.Seek(st.Size()-256*1024, 0)
			}
			b, e := io.ReadAll(io.LimitReader(f, 256*1024))
			f.Close()
			if e != nil {
				fmt.Println(e)
			} else {
				fmt.Println(string(b))
			}
		default:
			fmt.Println("Bitte1,2,3 oder0 wählen.")
		}
	}
}
