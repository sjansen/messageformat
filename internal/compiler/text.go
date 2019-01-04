package compiler

import (
	"strings"

	"github.com/sjansen/messageformat/ast"
	"golang.org/x/text/language"
)

type text struct {
	Value string
}

func newText(lang language.Tag, t *ast.Text) (*text, error) {
	return &text{Value: t.Value}, nil
}

func (t *text) format(b *strings.Builder, lang language.Tag, arguments map[string]interface{}) error {
	b.WriteString(t.Value)
	return nil
}
