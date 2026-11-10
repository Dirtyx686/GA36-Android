package main

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func put16(b []byte, o int, n uint16) { binary.LittleEndian.PutUint16(b[o:], n) }
func put32(b []byte, o int, n uint32) { binary.LittleEndian.PutUint32(b[o:], n) }
func manifestFixture(sdk int, split bool) []byte {
	words := []string{"manifest", "package", "com.example.game", "uses-sdk", "minSdkVersion", "split", "config.en"}
	pool := make([]byte, 28+4*len(words))
	put16(pool, 0, 1)
	put16(pool, 2, 28)
	put32(pool, 8, uint32(len(words)))
	put32(pool, 16, 0x100)
	put32(pool, 20, uint32(len(pool)))
	for i, s := range words {
		put32(pool, 28+i*4, uint32(len(pool)-28-4*len(words)))
		pool = append(pool, byte(len(s)), byte(len(s)))
		pool = append(pool, []byte(s)...)
		pool = append(pool, 0)
	}
	for len(pool)%4 != 0 {
		pool = append(pool, 0)
	}
	put32(pool, 4, uint32(len(pool)))
	node := func(tag uint32, attrs [][3]uint32) []byte {
		b := make([]byte, 36+20*len(attrs))
		put16(b, 0, 0x102)
		put16(b, 2, 16)
		put32(b, 4, uint32(len(b)))
		put32(b, 16, 0xffffffff)
		put32(b, 20, tag)
		put16(b, 24, 20)
		put16(b, 26, 20)
		put16(b, 28, uint16(len(attrs)))
		for i, a := range attrs {
			off := 36 + i*20
			put32(b, off, 0xffffffff)
			put32(b, off+4, a[0])
			put32(b, off+8, 0xffffffff)
			put16(b, off+12, 8)
			b[off+15] = byte(a[1])
			put32(b, off+16, a[2])
		}
		return b
	}
	attrs := [][3]uint32{{1, 3, 2}}
	if split {
		attrs = append(attrs, [3]uint32{5, 3, 6})
	}
	out := make([]byte, 8)
	put16(out, 0, 3)
	put16(out, 2, 8)
	out = append(out, pool...)
	out = append(out, node(0, attrs)...)
	out = append(out, node(3, [][3]uint32{{4, 0x10, uint32(sdk)}})...)
	put32(out, 4, uint32(len(out)))
	return out
}
func apkFixture(t *testing.T, sdk int, abi string, split bool) string {
	t.Helper()
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	w, _ := z.Create("AndroidManifest.xml")
	w.Write(manifestFixture(sdk, split))
	if abi != "" {
		w, _ = z.Create("lib/" + abi + "/libgame.so")
		w.Write([]byte("test ARM payload"))
	}
	z.Close()
	p := filepath.Join(t.TempDir(), "game.apk")
	os.WriteFile(p, b.Bytes(), 0600)
	return p
}
func inboxFixture(t *testing.T) string {
	t.Helper()
	p := t.TempDir()
	os.WriteFile(filepath.Join(p, "GA36-APPS.TXT"), []byte(inboxMarker), 0600)
	os.Mkdir(filepath.Join(p, "APKS"), 0700)
	os.Mkdir(filepath.Join(p, "OBB"), 0700)
	return p
}
func TestAPKCompatibility(t *testing.T) {
	for _, c := range []struct {
		sdk       int
		abi       string
		split, ok bool
	}{{24, "armeabi-v7a", false, true}, {21, "", false, true}, {25, "armeabi-v7a", false, false}, {24, "arm64-v8a", false, false}, {24, "x86", false, false}, {24, "armeabi-v7a", true, false}} {
		info, e := inspectAPK(apkFixture(t, c.sdk, c.abi, c.split))
		if (e == nil) != c.ok {
			t.Fatalf("%+v %v", c, e)
		}
		if e == nil && (info.Package != "com.example.game" || info.MinSDK != c.sdk) {
			t.Fatal(info)
		}
	}
}
func TestAPKQueueAtomicAndChecked(t *testing.T) {
	root := inboxFixture(t)
	info, e := queueAPK(root, apkFixture(t, 24, "armeabi-v7a", false))
	if e != nil {
		t.Fatal(e)
	}
	files, _ := filepath.Glob(filepath.Join(root, "APKS", "*.apk"))
	if len(files) != 1 {
		t.Fatal(files)
	}
	b, _ := os.ReadFile(files[0])
	if digest(b) != info.SHA {
		t.Fatal("copy mismatch")
	}
	pkg, _ := os.ReadFile(files[0] + ".pkg")
	if string(pkg) != info.Package+"\n" {
		t.Fatal("metadata mismatch")
	}
	if _, e = queueAPK(root, apkFixture(t, 26, "arm64-v8a", false)); e == nil {
		t.Fatal("accepted incompatible APK")
	}
	files, _ = filepath.Glob(filepath.Join(root, "APKS", "*.apk"))
	if len(files) != 1 {
		t.Fatal("incompatible queued")
	}
}
func TestOBBQueueAndExpansionVersion(t *testing.T) {
	root := inboxFixture(t)
	src := filepath.Join(t.TempDir(), "com.example.game")
	os.Mkdir(src, 0700)
	name := "main.11.com.example.game.obb"
	payload := bytes.Repeat([]byte("game-data"), 8192)
	os.WriteFile(filepath.Join(src, name), payload, 0600)
	pkg, e := queueOBB(root, src)
	if e != nil || pkg != "com.example.game" {
		t.Fatal(pkg, e)
	}
	jobs, _ := os.ReadDir(filepath.Join(root, "OBB"))
	if len(jobs) != 1 {
		t.Fatal(jobs)
	}
	job := filepath.Join(root, "OBB", jobs[0].Name())
	b, _ := os.ReadFile(filepath.Join(job, name))
	if !bytes.Equal(b, payload) {
		t.Fatal("OBB mismatch")
	}
	b, _ = os.ReadFile(filepath.Join(job, "SHA256.TXT"))
	if string(b) != digest(payload)+" "+name+"\n" {
		t.Fatal("bad OBB manifest")
	}
}
func TestOBBMalformedNameRejected(t *testing.T) {
	p := filepath.Join(t.TempDir(), "main.11.bad.obb")
	os.WriteFile(p, []byte("x"), 0600)
	if _, _, e := inspectOBB(p); e == nil {
		t.Fatal("accepted malformed package")
	}
}
func TestInboxMarkerRequired(t *testing.T) {
	if _, e := queueAPK(t.TempDir(), apkFixture(t, 24, "", false)); e == nil {
		t.Fatal("wrote outside approved inbox")
	}
}
func TestManifestTruncationNeverPanics(t *testing.T) {
	b := manifestFixture(24, false)
	for n := 0; n < len(b); n++ {
		if _, e := parseManifest(b[:n]); e == nil {
			t.Fatalf("truncated manifest accepted at%d", n)
		}
	}
	for i := range b {
		c := append([]byte(nil), b...)
		c[i] = 0xff
		parseManifest(c)
	}
}
func TestRealAPKOptional(t *testing.T) {
	p := os.Getenv("GA36_APK_TEST")
	if p == "" {
		t.Skip("no real APK path")
	}
	a, e := inspectAPK(p)
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(a.Package, ".") {
		t.Fatal(a)
	}
	t.Logf("Actual APK parsed: %s minSdk=%d ABI=%v", a.Package, a.MinSDK, a.Native)
}
