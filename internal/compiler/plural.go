package compiler

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"

	"github.com/sjansen/messageformat/ast"
)

type numberSign struct {
	ArgID  string
	Offset int
}

type pluralArg struct {
	ArgID    string
	Ordinal  bool
	Offset   int
	Messages map[string]*Message
}

func (n *numberSign) format(b *strings.Builder, arguments map[string]interface{}) error {
	value, ok := arguments[n.ArgID]
	if !ok {
		return fmt.Errorf("missing arg: %q", n.ArgID)
	}
	i, ok := value.(int)
	if !ok {
		return fmt.Errorf("expected int got: %T", value)
	}
	b.WriteString(strconv.Itoa(i - n.Offset))
	return nil
}

func newPluralArg(p *ast.PluralArg) (*pluralArg, error) {
	if _, ok := p.Messages["other"]; !ok {
		return nil, fmt.Errorf(`missing required plural category: "other"`)
	}
	n := &numberSign{
		ArgID:  p.ArgID,
		Offset: p.Offset,
	}
	messages := make(map[string]*Message, len(p.Messages))
	for k, v := range p.Messages {
		if msg, err := compile(v, n); err != nil {
			return nil, err
		} else {
			messages[k] = msg
		}
	}
	return &pluralArg{
		ArgID:    p.ArgID,
		Ordinal:  p.Ordinal,
		Offset:   p.Offset,
		Messages: messages,
	}, nil
}

func (p *pluralArg) format(b *strings.Builder, arguments map[string]interface{}) error {
	value, ok := arguments[p.ArgID]
	if !ok {
		return fmt.Errorf("missing arg: %q", p.ArgID)
	}

	n, ok := value.(int)
	if !ok {
		return fmt.Errorf("expected int got: %T", value)
	}

	category := fmt.Sprintf("=%d", n)
	if msg, ok := p.Messages[category]; ok {
		return msg.format(b, arguments)
	}

	var form plural.Form
	lang := language.MustParse("en")
	if p.Ordinal {
		form = plural.Ordinal.MatchPlural(lang, n, 0, 0, 0, 0)
	} else {
		form = plural.Cardinal.MatchPlural(lang, n, 0, 0, 0, 0)
	}

	category = "other"
	switch form {
	case plural.Zero:
		category = "zero"
	case plural.One:
		category = "one"
	case plural.Two:
		category = "two"
	case plural.Few:
		category = "few"
	case plural.Many:
		category = "many"
	}

	if msg, ok := p.Messages[category]; ok {
		return msg.format(b, arguments)
	}

	msg := p.Messages["other"]
	return msg.format(b, arguments)
}
