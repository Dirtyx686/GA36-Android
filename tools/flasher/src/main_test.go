package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func digest(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func fixture(t *testing.T) (*os.File, Manifest, []byte) {
	t.Helper()
	raw := make([]byte, 1024*1024+4096)
	for i := range raw {
		raw[i] = byte((i*37 + i/4096) % 251)
	}
	var b bytes.Buffer
	z := gzip.NewWriter(&b)
	z.Write(raw)
	z.Close()
	p := filepath.Join(t.TempDir(), "image.gz")
	os.WriteFile(p, b.Bytes(), 0600)
	f, e := os.Open(p)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { f.Close() })
	return f, Manifest{File: "image.gz", ImageBytes: int64(len(raw)), ImageSHA: digest(raw), HeaderSHA: digest(raw[:4096]), CompressedSHA: digest(b.Bytes())}, raw
}
func testDevice(t *testing.T) (*Device, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "simulated-sd.img")
	f, e := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		t.Fatal(e)
	}
	f.Write(bytes.Repeat([]byte{0xab}, 2*1024*1024))
	t.Cleanup(func() { f.Close() })
	return &Device{File: f, Guard: func() error { return nil }, ReopenRead: func() (*os.File, error) { return os.Open(path) }, Close: f.Close, Finish: f.Sync}, path
}
func TestActualWriteAndReadback(t *testing.T) {
	f, m, raw := fixture(t)
	d, p := testDevice(t)
	var phases []string
	e := flash(context.Background(), f, m, d, func(s string, a, b int64) { phases = append(phases, s) })
	if e != nil {
		t.Fatal(e)
	}
	got, _ := os.ReadFile(p)
	if !bytes.Equal(got[:len(raw)], raw) || got[len(raw)] != 0xab {
		t.Fatal("incorrect write bounds")
	}
	if len(phases) == 0 || phases[len(phases)-1] != "Zurücklesen und prüfen" {
		t.Fatal(phases)
	}
}
func TestCorruptImageNeverWrites(t *testing.T) {
	f, m, _ := fixture(t)
	m.CompressedSHA = digest([]byte("wrong"))
	d, p := testDevice(t)
	before, _ := os.ReadFile(p)
	if e := flash(context.Background(), f, m, d, func(string, int64, int64) {}); e == nil {
		t.Fatal("accepted corrupted image")
	}
	after, _ := os.ReadFile(p)
	if !bytes.Equal(before, after) {
		t.Fatal("wrote before image validation")
	}
}
func TestChangedCardNeverWrites(t *testing.T) {
	f, m, _ := fixture(t)
	d, p := testDevice(t)
	d.Guard = func() error { return errors.New("replaced") }
	before, _ := os.ReadFile(p)
	if flash(context.Background(), f, m, d, func(string, int64, int64) {}) == nil {
		t.Fatal("ignored identity guard")
	}
	after, _ := os.ReadFile(p)
	if !bytes.Equal(before, after) {
		t.Fatal("wrote after identity failure")
	}
}
func TestReadbackCorruptionDetected(t *testing.T) {
	f, m, _ := fixture(t)
	d, p := testDevice(t)
	d.ReopenRead = func() (*os.File, error) {
		os.WriteFile(p, bytes.Repeat([]byte{0}, int(m.ImageBytes)), 0600)
		return os.Open(p)
	}
	if flash(context.Background(), f, m, d, func(string, int64, int64) {}) == nil {
		t.Fatal("missed corruption")
	}
}
func TestWrongRawLengthAndHash(t *testing.T) {
	for _, which := range []string{"short", "long", "hash"} {
		t.Run(which, func(t *testing.T) {
			f, m, _ := fixture(t)
			d, _ := testDevice(t)
			switch which {
			case "short":
				m.ImageBytes += 4096
			case "long":
				m.ImageBytes -= 4096
			case "hash":
				m.ImageSHA = digest([]byte("wrong"))
			}
			if flash(context.Background(), f, m, d, func(string, int64, int64) {}) == nil {
				t.Fatal("accepted bad raw stream")
			}
		})
	}
}
func TestCancelled(t *testing.T) {
	f, m, _ := fixture(t)
	d, p := testDevice(t)
	before, _ := os.ReadFile(p)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if flash(ctx, f, m, d, func(string, int64, int64) {}) == nil {
		t.Fatal("ignored cancel")
	}
	after, _ := os.ReadFile(p)
	if !bytes.Equal(before, after) {
		t.Fatal("wrote after cancellation")
	}
}
func TestDiskFilter(t *testing.T) {
	m := Manifest{ImageBytes: 31138512896}
	for _, c := range []struct {
		d    Disk
		want bool
	}{{Disk{Safe: true, Bytes: m.ImageBytes}, true}, {Disk{Safe: true, Bytes: 64000000000}, true}, {Disk{Safe: true, Bytes: m.ImageBytes - 512}, false}, {Disk{Safe: false, Bytes: 64000000000}, false}, {Disk{Safe: true, Bytes: 3 * 1024 * 1024 * 1024 * 1024}, false}} {
		if eligible(c.d, m) != c.want {
			t.Fatal(c)
		}
	}
}
func TestManifestPathInjection(t *testing.T) {
	_, m, _ := fixture(t)
	for _, s := range []string{"../image.gz", "C:\\image.gz", "x/y.gz", "x\\y.gz"} {
		m.File = s
		if validManifest(m) == nil {
			t.Fatal("unsafe image name", s)
		}
	}
}
