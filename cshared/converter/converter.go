package converter

import "unsafe"

type Payload byte

const (
	PayloadTypeNull Payload = iota
	PlayloadTypeValue
)

type Converter struct {
}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) ConvertBytes(p Payload, body []byte) []byte {
	length := 0
	if body != nil {
		length = len(body)
	}
	return append([]byte{byte(p)}, (append(i32tob((uint32(length))), body...))...)
}

func (c *Converter) ConvertString(p Payload, body string) []byte {
	data := *(*[]byte)(unsafe.Pointer(&body))
	return c.ConvertBytes(p, data)
}

func i32tob(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}
