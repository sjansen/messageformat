package compiler

import "strings"

type Message struct {
	parts []part
}

type part interface {
	format(*strings.Builder, map[string]string) error
}

func (m *Message) Format(arguments map[string]string) (string, error) {
	var b strings.Builder
	if err := m.format(&b, arguments); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (m *Message) format(b *strings.Builder, arguments map[string]string) error {
	for _, part := range m.parts {
		if err := part.format(b, arguments); err != nil {
			return err
		}
	}
	return nil
}
