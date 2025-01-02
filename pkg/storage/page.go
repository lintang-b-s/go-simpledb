package storage

import (
	"bytes"
	"encoding/binary"
)

// Page . menyimpan data satu block page di dalam memori buffer (also disimpan di disk). (berukuran blockSize)
type Page struct {
	bb *bytes.Buffer
}

func NewPage(blockSize int) *Page {

	bb := bytes.NewBuffer(make([]byte, blockSize))
	return &Page{bb}
}

func NewPageFromByteSlice(b []byte) *Page {
	return &Page{bytes.NewBuffer(b)}
}



func (p *Page) GetInt(offset int) int {
	return int(binary.LittleEndian.Uint32(p.bb.Bytes()[offset:]))
}

// PutInt. set int ke byte array page di posisi = offset.
func (p *Page) PutInt(offset int, val int) {
	binary.LittleEndian.PutUint32(p.bb.Bytes()[offset:], uint32(val))
}

// GetBytes. return byte array dari byte array page di posisi = offset. di awal ada panjang bytes nya sehingga buat read bytes tinggal baca buffer page[offset+4:offset+4+length]
func (p *Page) GetBytes(offset int) []byte {
	length := p.GetInt(offset)
	b := make([]byte, p.GetInt(offset))
	copy(b, p.bb.Bytes()[offset+4:offset+4+length])
	return b
}

// PutBytes. set byte array ke byte array page di posisi = offset.
func (p *Page) PutBytes(offset int, b []byte) {
	p.PutInt(offset, len(b))
	copy(p.bb.Bytes()[offset+4:], b)
}

// GetString. return string dari byte array page di posisi= offset.
func (p *Page) GetString(offset int) string {
	return string(p.GetBytes(offset))
}

// putString. set string ke byte array page di posisi = offset.
func (p *Page) PutString(offset int, s string) {
	p.PutBytes(offset, []byte(s))
}

func (p *Page) Contents() []byte {
	return p.bb.Bytes()
}
