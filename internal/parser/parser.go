package parser

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/sjansen/messageformat/ast"
	"github.com/sjansen/messageformat/errors"
	"github.com/sjansen/messageformat/internal/lexer"
	"github.com/sjansen/messageformat/internal/lexer/token"
)

func Parse(r io.Reader) (*ast.Message, error) {
	l := lexer.New(r)
	nodes := []ast.Node{}
	for {
		tkn, err := l.Next()
		if err != nil {
			return nil, err
		}
		switch tkn.Type {
		case token.EOF:
			msg := &ast.Message{
				Nodes: nodes,
			}
			return msg, nil
		case token.LBRACE:
			if arg, err := parseArgument(l); err != nil {
				return nil, err
			} else {
				nodes = append(nodes, arg)
			}
		case token.TEXT:
			nodes = append(nodes, &ast.Text{Value: tkn.Value})
		default:
			err = &errors.UnexpectedToken{Token: tkn.Value}
			return nil, err
		}
	}
}

func parseArgument(l *lexer.Lexer) (ast.Node, error) {
	tkn, err := l.Next()
	if err != nil {
		return nil, err
	} else if tkn.Type != token.TEXT {
		err = &errors.UnexpectedToken{Token: tkn.Value}
		return nil, err
	}
	argNameOrNumber := tkn.Value

	tkn, err = l.Next()
	if err != nil {
		return nil, err
	} else if tkn.Type != token.RBRACE {
		err = &errors.UnexpectedToken{Token: tkn.Value}
		return nil, err
	}

	arg := &ast.PlainArg{ArgID: argNameOrNumber}
	return arg, nil
}

type parser struct {
	dec *Decoder
}

func (p *parser) parseArgument() (ast.Node, error) {
	p.dec.Decode()
	ch := p.dec.Decoded()
	if ch != '{' {
		return nil, &errors.UnexpectedToken{Token: string(ch)}
	}

	p.skipWhiteSpace()

	arg := &ast.PlainArg{ArgID: p.parseID()}

	p.skipWhiteSpace()

	p.dec.Decode()
	ch = p.dec.Decoded()
	if ch != '}' {
		return nil, &errors.UnexpectedToken{Token: string(ch)}
	}
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

func (p *parser) skipWhiteSpace() {
	for next := p.dec.Peek(); isPatternWhiteSpace(next); next = p.dec.Peek() {
		if !p.dec.Decode() {
			break
		}
	}
}
