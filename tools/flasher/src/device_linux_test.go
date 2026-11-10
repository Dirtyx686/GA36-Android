//go:build linux

package main

import (
	"context"
	"os"
	"syscall"
	"testing"
)

func TestSystemAndImageMountProtection(t *testing.T) {
	for _, m := range []string{"/", "/boot", "/boot/efi", "/home", "/usr", "/var"} {
		if !protectedMount(m, "/run/media/user/image/image.gz") {
			t.Fatal(m)
		}
	}
	if !protectedMount("/run/media/user/image", "/run/media/user/image/image.gz") {
		t.Fatal("image source unprotected")
	}
	if protectedMount("/run/media/user/target", "/run/media/user/image/image.gz") {
		t.Fatal("target unnecessarily rejected")
	}
}
func TestDirectReadbackOnTemporaryFile(t *testing.T) {
	f, m, _ := fixture(t)
	d, path := testDevice(t)
	d.ReopenRead = func() (*os.File, error) {
		if e := d.File.Sync(); e != nil {
			return nil, e
		}
		return os.OpenFile(path, os.O_RDONLY|syscall.O_DIRECT, 0)
	}
	if e := flash(context.Background(), f, m, d, func(string, int64, int64) {}); e != nil {
		t.Fatal(e)
	}
}
