package storage

import (
	"bytes"
	"encoding/binary"
)

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
	return int(p.bb.Bytes()[offset])
}

func (p *Page) PutInt(offset int, val int) {
	binary.LittleEndian.PutUint32(p.bb.Bytes()[offset:], uint32(val))
}

func (p *Page) GetBytes(offset int) []byte {
	length := p.GetInt(offset)
	b := make([]byte, p.GetInt(offset))
	copy(b, p.bb.Bytes()[offset+4:offset+4+length])
	return b
}

func (p *Page) PutBytes(offset int, b []byte) {
	p.PutInt(offset, len(b))
	copy(p.bb.Bytes()[offset+4:], b)
}


func (p *Page) GetString(offset int) string {
	return string(p.GetBytes(offset))
}

func (p *Page) putString(offset int, s string) {
	p.PutBytes(offset, []byte(s))
}

func (p *Page) contents() []byte {
	return p.bb.Bytes()
}