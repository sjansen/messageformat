package compiler

import (
	"strings"

	"golang.org/x/text/language"
)

type Message struct {
	lang  language.Tag
	parts []part
}

type part interface {
	format(*strings.Builder, language.Tag, map[string]interface{}) error
}

func (m *Message) Format(arguments map[string]interface{}) (string, error) {
	var b strings.Builder
	if err := m.format(&b, m.lang, arguments); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (m *Message) format(b *strings.Builder, lang language.Tag, arguments map[string]interface{}) error {
	for _, part := range m.parts {
		if err := part.format(b, lang, arguments); err != nil {
			return err
		}
	}
	return nil
}
