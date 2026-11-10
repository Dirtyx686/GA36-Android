package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
)

func headerFixture() []byte {
	b := make([]byte, 4096)
	for i := 512; i < len(b); i++ {
		b[i] = byte(i % 251)
	}
	b[510] = 0x55
	b[511] = 0xaa
	b[450] = 0x0c
	put32(b, 454, 3383336)
	put32(b, 458, 8388608)
	b[498] = 0x83
	put32(b, 502, uint32(dataOffset/512))
	put32(b, 506, uint32((baseImageBytes-dataOffset-endReserve)/512))
	return b
}
func TestCapacityPatchOnlyChangesDataLength(t *testing.T) {
	for _, cap := range []int64{baseImageBytes, 32000000000, 64000000000, 128000000000} {
		base := headerFixture()
		got, e := adjustedHeader(base, cap)
		if e != nil {
			t.Fatal(e)
		}
		wanted := (cap - dataOffset - endReserve) / 4096 * 4096
		if int64(binary.LittleEndian.Uint32(got[506:]))*512 != wanted {
			t.Fatal("wrong data size")
		}
		copy(got[506:510], base[506:510])
		if !bytes.Equal(base, got) {
			t.Fatal("modified unrelated boot/EBR data")
		}
	}
}
func TestCapacityPatchReadback(t *testing.T) {
	d, _ := testDevice(t)
	b := headerFixture()
	d.File.Seek(0, 0)
	d.File.Write(b)
	if e := configureDataPartition(d, 64000000000, digest(b)); e != nil {
		t.Fatal(e)
	}
	d.File.Seek(0, 0)
	got := make([]byte, 4096)
	io.ReadFull(d.File, got)
	want, _ := adjustedHeader(b, 64000000000)
	if !bytes.Equal(want, got) {
		t.Fatal("patch not committed")
	}
}
func TestCapacityPatchRejectsUnknownLayouts(t *testing.T) {
	for _, off := range []int{450, 454, 458, 498, 502, 506, 510, 511} {
		b := headerFixture()
		b[off] ^= 1
		if _, e := adjustedHeader(b, 64000000000); e == nil {
			t.Fatal("accepted changed layout", off)
		}
	}
	if _, e := adjustedHeader(headerFixture(), baseImageBytes-512); e == nil {
		t.Fatal("accepted small card")
	}
}
