package lexer

import (
	"bufio"
	"io"
	"strings"

	"github.com/sjansen/messageformat/ast/token"
)

type Lexer struct {
	buf        *bufio.Reader
	line       int
	byteOffset int
	runeOffset int
	inArgument bool
	next       struct {
		r    rune
		size int
		err  error
	}
}

func New(r io.Reader) *Lexer {
	l := &Lexer{}
	if buf, ok := r.(*bufio.Reader); ok {
		l.buf = buf
	} else {
		l.buf = bufio.NewReader(r)
	}
	return l
}

func (l *Lexer) Next() (*token.Token, error) {
	r, err := l.peek()
	if err != nil {
		return nil, err
	}
	switch {
	case !l.inArgument && r == '{':
		l.inArgument = true
		l.read()
		t := &token.Token{
			Type:  token.LBRACE,
			Value: "{",
		}
		return t, nil
	case l.inArgument && r == '}':
		l.inArgument = false
		l.read()
		t := &token.Token{
			Type:  token.RBRACE,
			Value: "}",
		}
		return t, nil
	case l.inArgument:
		return l.readArgToken()
	default:
		return l.readMessageToken()
	}
}

func (l *Lexer) peek() (rune, error) {
	if l.next.err == nil && l.next.size == 0 {
		r, size, err := l.buf.ReadRune()
		l.next.r = r
		l.next.size = size
		l.next.err = err
	}
	return l.next.r, l.next.err
}

func (l *Lexer) read() (rune, error) {
	err := l.next.err
	if err != nil {
		return 0, err
	}

	r := l.next.r
	size := l.next.size
	if size != 0 {
		l.next.r = 0
		l.next.size = 0
	} else {
		r, size, err = l.buf.ReadRune()
		if err != nil {
			return 0, err
		}
	}

	l.byteOffset += size
	l.runeOffset++
	if r == '\n' {
		l.line++
	}
	return r, nil
}

// argNameOrNumber = argName | argNumber
// argName = [^[[:Pattern_Syntax:][:Pattern_White_Space:]]]+
// argNumber = '0' | ('1'..'9' ('0'..'9')*)
func (l *Lexer) readArgNameOrNumber() (*token.Token, error) {
	var b strings.Builder
	for {
		r, err := l.peek()
		if err != nil && err != io.EOF {
			return nil, err
		} else if err == io.EOF || isPatternWhiteSpace(r) || isPatternSyntax(r) {
			t := &token.Token{
				Type:  token.TEXT,
				Value: b.String(),
			}
			return t, nil
		}
		l.read()
		b.WriteRune(r)
	}
}

func (l *Lexer) readArgToken() (*token.Token, error) {
	l.skipWhiteSpace()
	t, err := l.readArgNameOrNumber()
	l.skipWhiteSpace()
	return t, err
}

func (l *Lexer) readMessageToken() (*token.Token, error) {
	var b strings.Builder
	for {
		r, err := l.peek()
		if err != nil && err != io.EOF {
			return nil, err
		} else if err == io.EOF || r == '{' {
			t := &token.Token{
				Type:  token.TEXT,
				Value: b.String(),
			}
			return t, nil
		}
		l.read()
		b.WriteRune(r)
	}
}

func (l *Lexer) skipWhiteSpace() {
	for {
		if r, err := l.peek(); err != nil {
			return
		} else if !isPatternWhiteSpace(r) {
			return
		}
		l.read()
	}
}
