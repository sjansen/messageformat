package compiler

import (
	"fmt"

	"github.com/sjansen/messageformat/ast"
	"golang.org/x/text/language"
)

func Compile(lang string, msg *ast.Message) (*Message, error) {
	tag, err := language.Parse(lang)
	if err != nil {
		return nil, err
	}
	return compile(tag, msg, nil)
}

func compile(lang language.Tag, msg *ast.Message, n *numberSign) (*Message, error) {
	parts := make([]part, 0, len(msg.Parts))
	for _, part := range msg.Parts {
		switch x := part.(type) {
		case *ast.NumberSign:
			if n == nil {
				return nil, fmt.Errorf("illegal NumberSign")
			}
			parts = append(parts, n)
		case *ast.PlainArg:
			tmp, err := newPlainArg(lang, x)
			if err != nil {
				return nil, err
			}
			parts = append(parts, tmp)
		case *ast.PluralArg:
			tmp, err := newPluralArg(lang, x)
			if err != nil {
				return nil, err
			}
			parts = append(parts, tmp)
		case *ast.SelectArg:
			tmp, err := newSelectArg(lang, x)
			if err != nil {
				return nil, err
			}
			parts = append(parts, tmp)
		case *ast.Text:
			tmp, err := newText(lang, x)
			if err != nil {
				return nil, err
			}
			parts = append(parts, tmp)
		}
	}
	return &Message{lang: lang, parts: parts}, nil
}
