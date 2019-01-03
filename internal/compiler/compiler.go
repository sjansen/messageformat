package compiler

import (
	"fmt"

	"github.com/sjansen/messageformat/ast"
)

func Compile(msg *ast.Message) (*Message, error) {
	return compile(msg, nil)
}

func compile(msg *ast.Message, n *numberSign) (*Message, error) {
	parts := make([]part, 0, len(msg.Parts))
	for _, part := range msg.Parts {
		switch x := part.(type) {
		case *ast.NumberSign:
			if n == nil {
				return nil, fmt.Errorf("illegal NumberSign")
			} else {
				parts = append(parts, n)
			}
		case *ast.PlainArg:
			if tmp, err := newPlainArg(x); err != nil {
				return nil, err
			} else {
				parts = append(parts, tmp)
			}
		case *ast.PluralArg:
			if tmp, err := newPluralArg(x); err != nil {
				return nil, err
			} else {
				parts = append(parts, tmp)
			}
		case *ast.SelectArg:
			if tmp, err := newSelectArg(x); err != nil {
				return nil, err
			} else {
				parts = append(parts, tmp)
			}
		case *ast.Text:
			if tmp, err := newText(x); err != nil {
				return nil, err
			} else {
				parts = append(parts, tmp)
			}
		}
	}
	return &Message{parts: parts}, nil
}
