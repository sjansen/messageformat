package parser

import (
	"strings"
	"unicode/utf8"

	"github.com/sjansen/messageformat/ast"
	"github.com/sjansen/messageformat/errors"
)

func Parse(s string) (*ast.Message, error) {
	p := &parser{dec: NewDecoder(s)}
	if nodes, err := p.parseMessage(); err != nil {
		return nil, err
	} else {
		msg := &ast.Message{Nodes: nodes}
		return msg, nil
	}
}

type parser struct {
	dec *Decoder
}

func (p *parser) parseArgument() (ast.Node, error) {
	if err := p.requireRune('{'); err != nil {
		return nil, err
	}

	p.skipWhiteSpace()
	argNameOrNumber := p.parseID()
	p.skipWhiteSpace()

	p.dec.Decode()
	ch := p.dec.Decoded()
	if ch == '}' {
		arg := &ast.PlainArg{ArgID: argNameOrNumber}
		return arg, nil
	} else if ch != ',' {
		return nil, &errors.UnexpectedToken{Token: string(ch)}
	}

	p.skipWhiteSpace()
	keyword := p.parseID()
	argType := ast.ArgTypeFromKeyword(keyword)
	p.skipWhiteSpace()

	p.dec.Decode()
	ch = p.dec.Decoded()
	if ch == '}' {
		arg := &ast.SimpleArg{ArgID: argNameOrNumber, ArgType: argType}
		return arg, nil
	} else if ch != ',' {
		return nil, &errors.UnexpectedToken{Token: string(ch)}
	}

	p.skipWhiteSpace()
	keyword = p.parseID()
	argStyle := ast.ArgStyleFromKeyword(keyword)
	p.skipWhiteSpace()

	if err := p.requireRune('}'); err != nil {
		return nil, err
	}

	arg := &ast.SimpleArg{ArgID: argNameOrNumber, ArgType: argType, ArgStyle: argStyle}
	return arg, nil
}

func (p *parser) parseID() string {
	var b strings.Builder
	for p.dec.Decode() {
		ch := p.dec.Decoded()
		b.WriteRune(ch)
		next := p.dec.Peek()
		if isPatternWhiteSpace(next) || isPatternSyntax(next) {
			break
		}
	}
	return b.String()
}

func (p *parser) parseMessage() ([]ast.Node, error) {
	nodes := []ast.Node{}
	for {
		next := p.dec.Peek()
		if next == utf8.RuneError {
			break // TODO
		} else if next == '{' {
			node, err := p.parseArgument()
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		} else {
			node, err := p.parseMessageText()
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}

func (p *parser) parseMessageText() (*ast.Text, error) {
	inQuote := false
	var b strings.Builder
	for p.dec.Decode() {
		ch := p.dec.Decoded()
		if ch == '\'' {
			next := p.dec.Peek()
			if next == utf8.RuneError {
				if !inQuote {
					b.WriteRune('\'')
				}
				break
			} else if next == '\'' {
				b.WriteRune('\'')
				p.dec.Decode()
				next := p.dec.Peek()
				if !inQuote && next == '{' {
					break
				}
			} else if inQuote {
				inQuote = false
			} else if next == '{' || next == '}' {
				inQuote = true
			} else {
				b.WriteRune('\'')
			}
		} else {
			b.WriteRune(ch)
			next := p.dec.Peek()
			if next == '{' {
				break
			}
		}
	}
	t := &ast.Text{Value: b.String()}
	return t, nil
}

func (p *parser) requireRune(token rune) error {
	p.dec.Decode()
	ch := p.dec.Decoded()
	if ch == token {
		return nil
	}
	return &errors.UnexpectedToken{Token: string(ch)}
}

func (p *parser) skipWhiteSpace() {
	for next := p.dec.Peek(); isPatternWhiteSpace(next); next = p.dec.Peek() {
		if !p.dec.Decode() {
			break
		}
	}
}
