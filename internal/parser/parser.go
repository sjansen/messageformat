package parser

import (
	"io"
	"strings"

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

func (p *parser) parseMessageText() (*ast.Text, error) {
	var b strings.Builder
	for p.dec.Decode() {
		if ch := p.dec.Decoded(); ch == '{' {
			break
		} else {
			b.WriteRune(ch)
		}
	}
	t := &ast.Text{Value: b.String()}
	return t, nil
}
