package compiler

import (
	"fmt"
	"strings"

	"github.com/sjansen/messageformat/ast"
)

type selectArg struct {
	ArgID    string
	Messages map[string]*Message
}

func newSelectArg(s *ast.SelectArg) (*selectArg, error) {
	messages := make(map[string]*Message, len(s.Messages))
	for k, v := range s.Messages {
		if msg, err := Compile(v); err != nil {
			return nil, err
		} else {
			messages[k] = msg
		}
	}
	return &selectArg{
		ArgID:    s.ArgID,
		Messages: messages,
	}, nil
}

func (s *selectArg) format(b *strings.Builder, arguments map[string]interface{}) error {
	value, ok := arguments[s.ArgID]
	if !ok {
		return fmt.Errorf("missing arg: %q", s.ArgID)
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string got: %T", value)
	}
	msg, ok := s.Messages[str]
	if !ok {
		return fmt.Errorf("unmatched select: %q", value)
	}
	return msg.format(b, arguments)
}
