package lexer

import (
	"bufio"
	"io"
	"strings"

	"github.com/sjansen/messageformat/internal/lexer/token"
)

var EOF = &token.Token{
	Type:  token.EOF,
	Value: "",
}

type Lexer struct {
	buf        *bufio.Reader
	line       int
	byteOffset int
	runeOffset int

	inArgument bool

	next struct {
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

	c, size, err := l.buf.ReadRune()
	l.next.r = c
	l.next.size = size
	l.next.err = err
	return l
}

func (l *Lexer) Next() (*token.Token, error) {
	switch {
	case l.next.err == io.EOF:
		return EOF, nil
	case l.next.err != nil:
		return nil, l.next.err
	case l.inArgument:
		return l.readArgToken()
	default:
		return l.readNextToken()
	}
}

func (l *Lexer) peek() (rune, error) {
	return l.next.r, l.next.err
}

func (l *Lexer) read() (rune, error) {
	if l.next.err != nil {
		return 0, l.next.err
	}

	l.byteOffset += l.next.size
	l.runeOffset++
	if l.next.r == '\n' {
		l.line++
	}

	r, size, err := l.buf.ReadRune()
	r, l.next.r = l.next.r, r
	l.next.size = size
	l.next.err = err

	return r, nil
}

// argNameOrNumber = argName | argNumber
// argName = [^[[:Pattern_Syntax:][:Pattern_White_Space:]]]+
// argNumber = '0' | ('1'..'9' ('0'..'9')*)
func (l *Lexer) readArgNameOrNumber() (*token.Token, error) {
	s, err := l.readWhile(func(r rune) bool {
		return !isPatternWhiteSpace(r) && !isPatternSyntax(r)
	})
	if err != nil {
		return nil, err
	}

	t := &token.Token{
		Type:  token.TEXT,
		Value: s,
	}
	return t, nil
}

func (l *Lexer) readArgToken() (*token.Token, error) {
	l.skipWhiteSpace()

	r, err := l.peek()
	if err == nil {
		if r == ',' {
			_, _ = l.read()
			t := &token.Token{
				Type:  token.COMMA,
				Value: ",",
			}
			return t, nil
		} else if r == '}' {
			_, _ = l.read()
			l.inArgument = false
			t := &token.Token{
				Type:  token.RBRACE,
				Value: "}",
			}
			return t, nil
		}
	}

	return l.readArgNameOrNumber()
}

/*
messageText can contain quoted literal strings including syntax characters.
A quoted literal string begins with an ASCII apostrophe and a syntax
character (usually a {curly brace}) and continues until the next single
apostrophe. A double ASCII apostrohpe inside or outside of a quoted string
represents one literal apostrophe.

Quotable syntax characters are the {curly braces} in all messageText parts,
plus the '#' sign in a messageText immediately inside a pluralStyle,
and the '|' symbol in a messageText immediately inside a choiceStyle.
*/
func (l *Lexer) readMessageToken() (*token.Token, error) {
	var b strings.Builder
	for l.next.err == nil {
		if l.next.r == '{' {
			break
		} else if l.next.r == '\'' {
			_, _ = l.read()
			if l.next.r == '\'' {
				b.WriteRune('\'')
				_, _ = l.read()
				continue
			} else if l.next.r != '{' && l.next.r != '}' {
				b.WriteRune('\'')
				if l.next.err == nil {
					b.WriteRune(l.next.r)
					_, _ = l.read()
				}
				continue
			}
			for l.next.err == nil {
				if l.next.r == '\'' {
					_, _ = l.read()
					if l.next.r == '\'' {
						b.WriteRune('\'')
						_, _ = l.read()
					} else {
						break
					}
				} else {
					b.WriteRune(l.next.r)
					_, _ = l.read()
				}
			}
		} else {
			b.WriteRune(l.next.r)
			_, _ = l.read()
		}
	}

	t := &token.Token{
		Type:  token.TEXT,
		Value: b.String(),
	}
	return t, nil
}

func (l *Lexer) readNextToken() (*token.Token, error) {
	r, err := l.peek()
	if err == nil && r == '{' {
		_, _ = l.read()
		l.inArgument = true
		t := &token.Token{
			Type:  token.LBRACE,
			Value: "{",
		}
		return t, nil
	}

	return l.readMessageToken()
}

func (l *Lexer) readWhile(f func(r rune) bool) (string, error) {
	var b strings.Builder
	for l.next.err == nil {
		if f(l.next.r) {
			b.WriteRune(l.next.r)
			_, _ = l.read()
		} else {
			break
		}
	}
	return b.String(), nil
}

func (l *Lexer) skipWhiteSpace() {
	for l.next.err == nil {
		if isPatternWhiteSpace(l.next.r) {
			_, _ = l.read()
		} else {
			break
		}
	}
}
