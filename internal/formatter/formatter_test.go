package formatter

import (
	"testing"

	"github.com/sjansen/messageformat/ast"
	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	require := require.New(t)

	message := &ast.Message{Parts: []ast.Part{
		&ast.SelectArg{ArgID: "timespan",
			Messages: map[string]*ast.Message{
				"afternoon": &ast.Message{Parts: []ast.Part{
					&ast.Text{Value: "Boa tarde, "},
					&ast.PlainArg{ArgID: "name"},
					&ast.Text{Value: "."},
				}},
				"evening": &ast.Message{Parts: []ast.Part{
					&ast.Text{Value: "Boa noite, "},
					&ast.PlainArg{ArgID: "name"},
					&ast.Text{Value: "."},
				}},
				"other": &ast.Message{Parts: []ast.Part{
					&ast.Text{Value: "Bom dia, "},
					&ast.PlainArg{ArgID: "name"},
					&ast.Text{Value: "."},
				}},
			}},
	}}
	for _, tc := range []struct {
		expected  string
		arguments map[string]string
	}{
		{`Bom dia, Alice.`, map[string]string{
			"name":     "Alice",
			"timespan": "other",
		}},
		{`Boa tarde, Bob.`, map[string]string{
			"name":     "Bob",
			"timespan": "afternoon",
		}},
		{`Boa noite, Eve.`, map[string]string{
			"name":     "Eve",
			"timespan": "evening",
		}},
	} {
		actual, err := Format(message, tc.arguments)
		require.NoError(err)
		require.Equal(tc.expected, actual)
	}
}
