package compiler

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjansen/messageformat/ast"
)

var hello = &ast.Message{Parts: []ast.Part{
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

func TestCompileAndFormat(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		expected  string
		message   *ast.Message
		arguments map[string]interface{}
	}{{`Bom dia, Alice.`,
		hello, map[string]interface{}{
			"name":     "Alice",
			"timespan": "other",
		},
	}, {`Boa tarde, Bob.`,
		hello, map[string]interface{}{
			"name":     "Bob",
			"timespan": "afternoon",
		},
	}, {`Boa noite, Eve.`,
		hello, map[string]interface{}{
			"name":     "Eve",
			"timespan": "evening",
		}},
	} {
		compiled, err := Compile(tc.message)
		require.NoError(err)

		actual, err := compiled.Format(tc.arguments)
		require.NoError(err)
		require.Equal(tc.expected, actual)
	}
}
