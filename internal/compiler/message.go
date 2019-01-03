package compiler

import "strings"

type Message struct {
	parts []part
}

type part interface {
	format(*strings.Builder, map[string]interface{}) error
}

func (m *Message) Format(arguments map[string]interface{}) (string, error) {
	var b strings.Builder
	if err := m.format(&b, arguments); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (m *Message) format(b *strings.Builder, arguments map[string]interface{}) error {
	for _, part := range m.parts {
		if err := part.format(b, arguments); err != nil {
			return err
		}
	}
	return nil
}
