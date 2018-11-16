package lexer

import (
	"strconv"
	"strings"
	"testing"

	"github.com/sjansen/messageformat/ast/token"
	"github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	for idx, tc := range []struct {
		input    string
		expected []*token.Token
	}{
		{"Spoon!", []*token.Token{{
			Type:  token.TEXT,
			Value: "Spoon!",
		}}},
		{"Olá mundo!", []*token.Token{{
			Type:  token.TEXT,
			Value: "Olá mundo!",
		}}},
		{"Hello, {audience}!", []*token.Token{{
			Type:  token.TEXT,
			Value: "Hello, ",
		}, {
			Type:  token.LBRACE,
			Value: "{",
		}, {
			Type:  token.TEXT,
			Value: "audience",
		}, {
			Type:  token.RBRACE,
			Value: "}",
		}, {
			Type:  token.TEXT,
			Value: "!",
		}}},
		{"{ greeting }, World!", []*token.Token{{
			Type:  token.LBRACE,
			Value: "{",
		}, {
			Type:  token.TEXT,
			Value: "greeting",
		}, {
			Type:  token.RBRACE,
			Value: "}",
		}, {
			Type:  token.TEXT,
			Value: ", World!",
		}}},
		{"It's peanut butter jelly time!", []*token.Token{{
			Type:  token.TEXT,
			Value: "It's peanut butter jelly time!",
		}}},
		{"It''s peanut butter jelly time!", []*token.Token{{
			Type:  token.TEXT,
			Value: "It's peanut butter jelly time!",
		}}},
	} {
		tc := tc
		name := strconv.Itoa(idx)
		t.Run(name, func(t *testing.T) {
			require := require.New(t)

			r := strings.NewReader(tc.input)
			l := New(r)
			for _, expected := range tc.expected {
				actual, err := l.Next()
				require.NoError(err)
				require.Equal(expected, actual)
			}
			tkn, err := l.Next()
			require.NoError(err)
			require.NotNil(tkn)
			require.Equal(token.EOF, tkn.Type)
		})
	}
}
