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
			} else {
				parts = append(parts, n)
			}
		case *ast.PlainArg:
			if tmp, err := newPlainArg(lang, x); err != nil {
				return nil, err
			} else {
				parts = append(parts, tmp)
			}
		case *ast.PluralArg:
			if tmp, err := newPluralArg(lang, x); err != nil {
				return nil, err
			} else {
				parts = append(parts, tmp)
			}
		case *ast.SelectArg:
			if tmp, err := newSelectArg(lang, x); err != nil {
				return nil, err
			} else {
				parts = append(parts, tmp)
			}
		case *ast.Text:
			if tmp, err := newText(lang, x); err != nil {
				return nil, err
			} else {
				parts = append(parts, tmp)
			}
		}
	}
	return &Message{lang: lang, parts: parts}, nil
}
