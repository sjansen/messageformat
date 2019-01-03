package compiler

import (
	"fmt"
	"strconv"
	"strings"

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

	category := "other"
	switch n {
	case 1:
		category = "one"
	case 2:
		category = "two"
	case 3:
		category = "few"
	}

	msg, ok := p.Messages[category]
	if !ok {
		msg = p.Messages["other"]
	}

	return msg.format(b, arguments)
}
