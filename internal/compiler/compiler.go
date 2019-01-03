package compiler

import "github.com/sjansen/messageformat/ast"

func Compile(msg *ast.Message) (*Message, error) {
	parts := make([]part, 0, len(msg.Parts))
	for _, part := range msg.Parts {
		switch x := part.(type) {
		case *ast.PlainArg:
			if tmp, err := newPlainArg(x); err != nil {
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
