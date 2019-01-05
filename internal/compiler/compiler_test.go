package compiler

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjansen/messageformat/ast"
)

var hello = &ast.Message{Parts: []ast.Part{
	&ast.SelectArg{ArgID: "timespan",
		Messages: map[string]*ast.Message{
			"afternoon": {Parts: []ast.Part{
				&ast.Text{Value: "Boa tarde, "},
				&ast.PlainArg{ArgID: "name"},
				&ast.Text{Value: "."},
			}},
			"evening": {Parts: []ast.Part{
				&ast.Text{Value: "Boa noite, "},
				&ast.PlainArg{ArgID: "name"},
				&ast.Text{Value: "."},
			}},
			"other": {Parts: []ast.Part{
				&ast.Text{Value: "Bom dia, "},
				&ast.PlainArg{ArgID: "name"},
				&ast.Text{Value: "."},
			}},
		}},
}}

var items = &ast.Message{Parts: []ast.Part{
	&ast.PluralArg{
		ArgID:   "count",
		Ordinal: true,
		Messages: map[string]*ast.Message{
			"one": {Parts: []ast.Part{
				&ast.NumberSign{},
				&ast.Text{Value: "st item"},
			}},
			"two": {Parts: []ast.Part{
				&ast.NumberSign{},
				&ast.Text{Value: "nd item"},
			}},
			"few": {Parts: []ast.Part{
				&ast.NumberSign{},
				&ast.Text{Value: "rd item"},
			}},
			"other": {Parts: []ast.Part{
				&ast.NumberSign{},
				&ast.Text{Value: "th item"},
			}},
		}},
}}

var elves = &ast.Message{Parts: []ast.Part{
	&ast.PluralArg{
		ArgID: "count",
		Messages: map[string]*ast.Message{
			"=0":    {Parts: []ast.Part{&ast.Text{Value: "no elves"}}},
			"one":   {Parts: []ast.Part{&ast.Text{Value: "one elf"}}},
			"other": {Parts: []ast.Part{&ast.Text{Value: "multiple elves"}}},
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
		},
	}, {`1st item`,
		items, map[string]interface{}{
			"count": 1,
		},
	}, {`2nd item`,
		items, map[string]interface{}{
			"count": 2,
		},
	}, {`3rd item`,
		items, map[string]interface{}{
			"count": 3,
		},
	}, {`4th item`,
		items, map[string]interface{}{
			"count": 4,
		},
	}, {`no elves`,
		elves, map[string]interface{}{
			"count": 0,
		},
	}, {`one elf`,
		elves, map[string]interface{}{
			"count": 1,
		},
	}, {`multiple elves`,
		elves, map[string]interface{}{
			"count": 2,
		},
	}} {
		compiled, err := Compile("en", tc.message)
		require.NoError(err)

		actual, err := compiled.Format(tc.arguments)
		require.NoError(err)
		require.Equal(tc.expected, actual)
	}
}
