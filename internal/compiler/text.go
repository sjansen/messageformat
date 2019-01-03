package compiler

import (
	"strings"

	"github.com/sjansen/messageformat/ast"
)

type text struct {
	Value string
}

func newText(t *ast.Text) (*text, error) {
	return &text{Value: t.Value}, nil
}

func (t *text) format(b *strings.Builder, arguments map[string]interface{}) error {
	b.WriteString(t.Value)
	return nil
}
