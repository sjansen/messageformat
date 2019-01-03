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

func (p *plainArg) format(b *strings.Builder, arguments map[string]interface{}) error {
	value, ok := arguments[p.ArgID]
	if !ok {
		return fmt.Errorf("missing arg: %q", p.ArgID)
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string got: %T", value)
	}
	b.WriteString(str)
	return nil
}
