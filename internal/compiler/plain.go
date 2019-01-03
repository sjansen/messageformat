package compiler

import (
	"fmt"
	"strings"

	"github.com/sjansen/messageformat/ast"
)

type plainArg struct {
	ArgID string
}

func newPlainArg(p *ast.PlainArg) (*plainArg, error) {
	return &plainArg{ArgID: p.ArgID}, nil
}

func (p *plainArg) format(b *strings.Builder, arguments map[string]string) error {
	value, ok := arguments[p.ArgID]
	if !ok {
		return fmt.Errorf("missing arg: %q", p.ArgID)
	}
	b.WriteString(value)
	return nil
}
