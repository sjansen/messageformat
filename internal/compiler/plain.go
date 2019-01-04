package compiler

import (
	"fmt"
	"strings"

	"github.com/sjansen/messageformat/ast"
	"golang.org/x/text/language"
)

type plainArg struct {
	ArgID string
}

func newPlainArg(lang language.Tag, p *ast.PlainArg) (*plainArg, error) {
	return &plainArg{ArgID: p.ArgID}, nil
}

func (p *plainArg) format(b *strings.Builder, lang language.Tag, arguments map[string]interface{}) error {
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
