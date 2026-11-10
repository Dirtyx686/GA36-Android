package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
)

const dataOffset = int64(8589934592)
const endReserve = int64(16777216)
const baseImageBytes = int64(31138512896)

func adjustedHeader(b []byte, capacity int64) ([]byte, error) {
	if len(b) != 4096 || capacity < baseImageBytes || capacity > 2*1024*1024*1024*1024 || capacity%512 != 0 {
		return nil, errors.New("Ungültige Kartengröße")
	}
	get := func(o int) uint32 { return binary.LittleEndian.Uint32(b[o:]) }
	if b[510] != 0x55 || b[511] != 0xaa || b[450] != 0x0c || get(454) != 3383336 || get(458) != 8388608 || b[498] != 0x83 || get(502) != uint32(dataOffset/512) || get(506) != uint32((baseImageBytes-dataOffset-endReserve)/512) {
		return nil, errors.New("Unerwartete Partitionstabelle; keine Änderung vorgenommen")
	}
	n := (capacity - dataOffset - endReserve) / 4096 * 4096
	out := append([]byte(nil), b...)
	binary.LittleEndian.PutUint32(out[506:], uint32(n/512))
	return out, nil
}
func configureDataPartition(d *Device, capacity int64, expectedHeaderSHA string) error {
	if e := d.Guard(); e != nil {
		return e
	}
	buf := alignedBuffer(4096)
	if _, e := d.File.Seek(0, 0); e != nil {
		return e
	}
	if _, e := io.ReadFull(d.File, buf); e != nil {
		return e
	}
	sum := sha256.Sum256(buf)
	if hex.EncodeToString(sum[:]) != expectedHeaderSHA {
		return errors.New("Boot-Metadaten haben sich geändert; keine Größenanpassung vorgenommen")
	}
	wanted, e := adjustedHeader(buf, capacity)
	if e != nil {
		return e
	}
	copy(buf, wanted)
	if _, e = d.File.Seek(0, 0); e != nil {
		return e
	}
	if e = writeFull(d.File, buf); e != nil {
		return e
	}
	if e = d.File.Sync(); e != nil {
		return e
	}
	r, e := d.ReopenRead()
	if e != nil {
		return e
	}
	defer r.Close()
	if _, e = io.ReadFull(r, buf); e != nil {
		return e
	}
	if !bytes.Equal(buf, wanted) {
		return errors.New("Größenanpassung konnte nicht zurückgeprüft werden")
	}
	return d.Guard()
}
