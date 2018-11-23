package parser

import (
	"fmt"
	"io"

	"github.com/sjansen/messageformat/ast"
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
			err = fmt.Errorf("Unexpected token: %q", tkn.Value)
			return nil, err
		}
	}
}

func parseArgument(l *lexer.Lexer) (ast.Node, error) {
	tkn, err := l.Next()
	if err != nil {
		return nil, err
	} else if tkn.Type != token.TEXT {
		err = fmt.Errorf("Unexpected token: %q", tkn.Value)
		return nil, err
	}
	argNameOrNumber := tkn.Value

	tkn, err = l.Next()
	if err != nil {
		return nil, err
	} else if tkn.Type != token.RBRACE {
		err = fmt.Errorf("Unexpected token: %q", tkn.Value)
		return nil, err
	}

	arg := &ast.PlainArg{ArgID: argNameOrNumber}
	return arg, nil
}
