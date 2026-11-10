// Copyright 2026 GA36 project. Original flasher code: MIT.
package main

import (
	"bufio"
	"compress/gzip"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

//go:embed release.json
var releaseJSON []byte

type Manifest struct {
	Version       string `json:"version"`
	File          string `json:"file"`
	CompressedSHA string `json:"compressed_sha256"`
	ImageSHA      string `json:"image_sha256"`
	HeaderSHA     string `json:"header_sha256"`
	ImageBytes    int64  `json:"image_bytes"`
	Board         string `json:"board"`
}
type Disk struct {
	Path, Model, Serial, ID, Bus, Identity, Seq string
	Bytes                                       int64
	Safe                                        bool
	Raw                                         bool
	Mounts                                      []string
	Volumes                                     []string
}
type Device struct {
	File       *os.File
	Guard      func() error
	ReopenRead func() (*os.File, error)
	Close      func() error
	Finish     func() error
}

func validManifest(m Manifest) error {
	if m.ImageBytes < 4096 || m.ImageBytes%4096 != 0 || m.File != filepath.Base(m.File) || strings.ContainsAny(m.File, "/\\") {
		return errors.New("Ungültiges Image-Paket")
	}
	for _, h := range []string{m.ImageSHA, m.CompressedSHA, m.HeaderSHA} {
		b, e := hex.DecodeString(h)
		if e != nil || len(b) != 32 {
			return errors.New("Image-Paket ist noch nicht freigegeben")
		}
	}
	return nil
}
func eligible(d Disk, m Manifest) bool {
	return d.Safe && d.Bytes >= m.ImageBytes && d.Bytes <= 2*1024*1024*1024*1024
}
func hashFile(f *os.File) (string, error) {
	if _, e := f.Seek(0, 0); e != nil {
		return "", e
	}
	h := sha256.New()
	if _, e := io.Copy(h, f); e != nil {
		return "", e
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
func validateImage(f *os.File, m Manifest) error {
	if e := validManifest(m); e != nil {
		return e
	}
	h, e := hashFile(f)
	if e != nil {
		return e
	}
	if h != m.CompressedSHA {
		return errors.New("Image-Prüfsumme stimmt nicht. Nichts wurde geschrieben.")
	}
	_, e = f.Seek(0, 0)
	return e
}

// A page-aligned buffer supports direct disk I/O on Linux and Windows.
func alignedBuffer(n int) []byte {
	b := make([]byte, n+4096)
	off := int((4096 - uintptr(unsafe.Pointer(&b[0]))%4096) % 4096)
	return b[off : off+n]
}
func writeFull(w io.Writer, b []byte) error {
	for len(b) > 0 {
		n, e := w.Write(b)
		if e != nil {
			return e
		}
		if n <= 0 || n > len(b) {
			return io.ErrShortWrite
		}
		b = b[n:]
	}
	return nil
}
func flash(ctx context.Context, img *os.File, m Manifest, d *Device, progress func(string, int64, int64)) error {
	if e := validateImage(img, m); e != nil {
		return e
	}
	if e := d.Guard(); e != nil {
		return e
	}
	gz, e := gzip.NewReader(img)
	if e != nil {
		return e
	}
	defer gz.Close()
	gz.Multistream(false)
	if _, e = d.File.Seek(0, 0); e != nil {
		return e
	}
	buf := alignedBuffer(1024 * 1024)
	h := sha256.New()
	var at int64
	for at < m.ImageBytes {
		if e = ctx.Err(); e != nil {
			return e
		}
		if at%(64*1024*1024) == 0 {
			if e = d.Guard(); e != nil {
				return e
			}
		}
		n := int64(len(buf))
		if n > m.ImageBytes-at {
			n = m.ImageBytes - at
		}
		if _, e = io.ReadFull(gz, buf[:n]); e != nil {
			return fmt.Errorf("Image endet zu früh: %w", e)
		}
		if e = writeFull(d.File, buf[:n]); e != nil {
			return fmt.Errorf("Schreibfehler: %w", e)
		}
		h.Write(buf[:n])
		at += n
		progress("Schreiben", at, m.ImageBytes)
	}
	var extra [1]byte
	n, e := gz.Read(extra[:])
	if n != 0 || e != io.EOF {
		return errors.New("Image-Länge oder GZIP-Abschluss ungültig")
	}
	if hex.EncodeToString(h.Sum(nil)) != m.ImageSHA {
		return errors.New("Entpackte Image-Prüfsumme stimmt nicht")
	}
	if e = d.File.Sync(); e != nil {
		return e
	}
	if e = d.Guard(); e != nil {
		return e
	}
	r, e := d.ReopenRead()
	if e != nil {
		return e
	}
	defer r.Close()
	h = sha256.New()
	at = 0
	for at < m.ImageBytes {
		if e = ctx.Err(); e != nil {
			return e
		}
		if at%(64*1024*1024) == 0 {
			if e = d.Guard(); e != nil {
				return e
			}
		}
		n := int64(len(buf))
		if n > m.ImageBytes-at {
			n = m.ImageBytes - at
		}
		if _, e = io.ReadFull(r, buf[:n]); e != nil {
			return fmt.Errorf("Lesefehler bei der Prüfung: %w", e)
		}
		h.Write(buf[:n])
		at += n
		progress("Zurücklesen und prüfen", at, m.ImageBytes)
	}
	if hex.EncodeToString(h.Sum(nil)) != m.ImageSHA {
		return errors.New("Schreibprüfung fehlgeschlagen. Karte nicht verwenden.")
	}
	return d.Guard()
}
func ask(r *bufio.Reader, p string) (string, error) {
	fmt.Print(p)
	s, e := r.ReadString('\n')
	if e != nil {
		return "", e
	}
	return strings.TrimSpace(s), nil
}
func run() error {
	var m Manifest
	if e := json.Unmarshal(releaseJSON, &m); e != nil {
		return e
	}
	args := os.Args[1:]
	checkOnly := len(args) == 1 && args[0] == "--check"
	versionOnly := len(args) == 1 && args[0] == "--version"
	if len(args) > 0 && !checkOnly && !versionOnly && !(len(args) == 1 && args[0] == "--elevated") {
		return errors.New("Unbekannte Option. Programm ohne Optionen starten.")
	}
	if versionOnly {
		fmt.Printf("GA36 Flasher %s (%s)\n", m.Version, m.Board)
		return nil
	}
	exe, e := os.Executable()
	if e != nil {
		return e
	}
	exe, e = filepath.EvalSymlinks(exe)
	if e != nil {
		return e
	}
	path := filepath.Join(filepath.Dir(exe), m.File)
	fmt.Printf("\nGA36 · Android 7 · %s\nNur für %s. Nicht für andere R36S-Boards.\n\n", m.Version, m.Board)
	in := bufio.NewReader(os.Stdin)
	if !checkOnly {
		if e = ensureAdmin(); e != nil {
			return e
		}
		fmt.Println("1 Android7 auf SD-Karte schreiben (löscht Zielkarte)\n2 APKs / OBB-Spieldaten auf bestehender Karte vorbereiten\n0 Ende")
		menu := in
		choice, err := ask(menu, "Auswahl: ")
		if err != nil {
			return err
		}
		if choice == "0" || choice == "" {
			return nil
		}
		if choice == "2" {
			return importMenu(exe, in)
		}
		if choice != "1" {
			return errors.New("Ungültige Auswahl")
		}
	}
	img, e := os.Open(path)
	if e != nil {
		return fmt.Errorf("Image neben dem Flasher fehlt. Paket vollständig entpacken: %w", e)
	}
	defer img.Close()
	fmt.Println("Image-Paket wird geprüft …")
	if e = validateImage(img, m); e != nil {
		return e
	}
	fmt.Println("Image-Prüfsumme: OK")
	if checkOnly {
		fmt.Println("Entpacktes Image wird vollständig geprüft …")
		gz, err := gzip.NewReader(img)
		if err != nil {
			return err
		}
		h := sha256.New()
		n, err := io.Copy(h, gz)
		closeErr := gz.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
		if n != m.ImageBytes || hex.EncodeToString(h.Sum(nil)) != m.ImageSHA {
			return errors.New("Entpackte Image-Länge oder Prüfsumme stimmt nicht")
		}
		fmt.Println("Vollständiges Image: OK")
		return nil
	}
	if e = ensureAdmin(); e != nil {
		return e
	}
	disks, e := listDisks(path, exe)
	if e != nil {
		return e
	}
	var choices []Disk
	for _, d := range disks {
		if eligible(d, m) {
			choices = append(choices, d)
		}
	}
	if len(choices) == 0 {
		return fmt.Errorf("Keine geeignete SD-Karte gefunden. Mindestens %.2f GB erforderlich; interne Systemlaufwerke sind gesperrt.", float64(m.ImageBytes)/1e9)
	}
	fmt.Println("\nKartenleser / externe Ziele:")
	for i, d := range choices {
		fmt.Printf("  %d. %s · %.2f GB · %s · %s\n", i+1, d.Model, float64(d.Bytes)/1e9, d.Path, d.Serial)
	}
	// Reuse the menu reader so piped/quick input is not lost.
	s, e := ask(in, "\nNummer der Zielkarte (Enter = Abbruch): ")
	if e != nil {
		return e
	}
	if s == "" {
		return errors.New("Abgebrochen; nichts geschrieben")
	}
	i, e := strconv.Atoi(s)
	if e != nil || i < 1 || i > len(choices) {
		return errors.New("Ungültige Auswahl; nichts geschrieben")
	}
	selected := choices[i-1]
	fmt.Printf("\nACHTUNG: Dateien und Partitionen auf %s (%s, %.2f GB) werden gelöscht!\nDie Konsole muss ausgeschaltet sein. Karte im PC-Kartenleser verwenden.\n", selected.Path, selected.Model, float64(selected.Bytes)/1e9)
	expected := fmt.Sprintf("LOESCHEN %d", i)
	s, e = ask(in, "Zum Bestätigen genau "+expected+" eingeben: ")
	if e != nil || s != expected {
		return errors.New("Abgebrochen; nichts geschrieben")
	}
	// Re-enumerate and compare the target before unmounting/opening it.
	current, e := listDisks(path, exe)
	if e != nil {
		return e
	}
	var match *Disk
	for j := range current {
		d := &current[j]
		if d.Path == selected.Path && d.Identity == selected.Identity && eligible(*d, m) {
			match = d
			break
		}
	}
	if match == nil {
		return errors.New("Zielkarte hat sich geändert. Bitte neu starten.")
	}
	dev, e := openDevice(*match, path, exe)
	if e != nil {
		return e
	}
	defer dev.Close()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	last := time.Time{}
	progress := func(phase string, at, total int64) {
		if time.Since(last) >= time.Second || at == total {
			fmt.Printf("\r%-25s %6.2f %%", phase, 100*float64(at)/float64(total))
			last = time.Now()
		}
	}
	if e = flash(ctx, img, m, dev, progress); e != nil {
		return fmt.Errorf("\nNicht abgeschlossen: %w\nDie Zielkarte kann unvollständig sein. Vollständig erneut flashen.", e)
	}
	fmt.Println("\nDatenpartition an Kartengröße anpassen …")
	if e = configureDataPartition(dev, selected.Bytes, m.HeaderSHA); e != nil {
		return e
	}
	if e = dev.Finish(); e != nil {
		return e
	}
	if e = dev.Close(); e != nil {
		return e
	}
	fmt.Println("\n\nFertig. Image vollständig geschrieben und zurückgelesen.\nKarte sicher entnehmen und in den Systemslot der GA36 stecken.\nBeim ersten Start wird der freie Platz automatisch eingerichtet. Nicht ausschalten.\nDer erste Start kann länger dauern; das Image wartet noch auf den Gerätetest.")
	return nil
}
func main() {
	args := os.Args[1:]
	if len(args) == 0 || (len(args) > 0 && strings.HasPrefix(args[0], "--gui")) {
		if e := runGUI(args); e != nil {
			fmt.Fprintln(os.Stderr, e)
			guiShowError(e)
			os.Exit(1)
		}
		return
	}
	if len(args) > 0 && args[0] == "--console" {
		os.Args = append([]string{os.Args[0]}, args[1:]...)
	}

	code := 0
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, "\n"+e.Error())
		code = 1
	}
	pauseExit()
	os.Exit(code)
}
