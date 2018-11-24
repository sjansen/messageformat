package parser

import "unicode/utf8"

type Decoder struct {
	src string
	idx int

	currRune rune
	currSize int
	nextRune rune
	nextSize int
}

func NewDecoder(s string) *Decoder {
	ch, size := utf8.DecodeRuneInString(s)
	return &Decoder{
		src:      s,
		idx:      0,
		currRune: utf8.RuneError,
		currSize: 0,
		nextRune: ch,
		nextSize: size,
	}
}

func (d *Decoder) Decode() bool {
	if d.nextSize < 1 {
		return false
	}

	d.currRune = d.nextRune
	d.currSize = d.nextSize
	d.idx += d.currSize

	ch, size := utf8.DecodeRuneInString(d.src[d.idx:])
	d.nextRune = ch
	d.nextSize = size

	return true
}

func (d *Decoder) Decoded() rune {
	return d.currRune
}

func (d *Decoder) Peek() rune {
	return d.nextRune
}
