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

func (s *selectArg) format(b *strings.Builder, arguments map[string]string) error {
	value, ok := arguments[s.ArgID]
	if !ok {
		return fmt.Errorf("missing arg: %q", s.ArgID)
	}
	msg, ok := s.Messages[value]
	if !ok {
		return fmt.Errorf("unmatched select: %q", value)
	}
	return msg.format(b, arguments)
}
