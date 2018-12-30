package formatter

import (
	"fmt"
	"strings"

	"github.com/sjansen/messageformat/ast"
)

func Format(msg *ast.Message, arguments map[string]string) (string, error) {
	var b strings.Builder
	if err := format(&b, msg.Parts, arguments); err != nil {
		return "", err
	}
	return b.String(), nil
}

func format(b *strings.Builder, parts []ast.Part, arguments map[string]string) error {
	for _, part := range parts {
		switch x := part.(type) {
		case *ast.PlainArg:
			if err := formatPlainArg(b, x, arguments); err != nil {
				return err
			}
		case *ast.SelectArg:
			if err := formatSelectArg(b, x, arguments); err != nil {
				return err
			}
		case *ast.Text:
			b.WriteString(x.Value)
		}
	}
	return nil
}

func formatPlainArg(b *strings.Builder, arg *ast.PlainArg, arguments map[string]string) error {
	value, ok := arguments[arg.ArgID]
	if !ok {
		return fmt.Errorf("missing arg: %q", arg.ArgID)
	}
	b.WriteString(value)
	return nil
}

func formatSelectArg(b *strings.Builder, arg *ast.SelectArg, arguments map[string]string) error {
	value, ok := arguments[arg.ArgID]
	if !ok {
		return fmt.Errorf("missing arg: %q", arg.ArgID)
	}
	msg, ok := arg.Messages[value]
	if !ok {
		return fmt.Errorf("unmatched select: %q", value)
	}
	if err := format(b, msg.Parts, arguments); err != nil {
		return err
	}
	return nil
}
