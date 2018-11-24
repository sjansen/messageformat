package parser

import (
	"strconv"
	"strings"
	"testing"

	"github.com/sjansen/messageformat/ast"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	for idx, tc := range []struct {
		input    string
		expected *ast.Message
	}{
		{"Spoon!", &ast.Message{Nodes: []ast.Node{
			&ast.Text{Value: "Spoon!"},
		}}},
		{"Olá mundo!", &ast.Message{Nodes: []ast.Node{
			&ast.Text{Value: "Olá mundo!"},
		}}},
		{"Hello, {audience}!", &ast.Message{Nodes: []ast.Node{
			&ast.Text{Value: "Hello, "},
			&ast.PlainArg{ArgID: "audience"},
			&ast.Text{Value: "!"},
		}}},
		{"{ greeting }, World!", &ast.Message{Nodes: []ast.Node{
			&ast.PlainArg{ArgID: "greeting"},
			&ast.Text{Value: ", World!"},
		}}},
		{"It's peanut butter jelly time!", &ast.Message{Nodes: []ast.Node{
			&ast.Text{Value: "It's peanut butter jelly time!"},
		}}},
		{"It''s peanut butter jelly time!", &ast.Message{Nodes: []ast.Node{
			&ast.Text{Value: "It's peanut butter jelly time!"},
		}}},
		{"'-'''{-''-}'''-'", &ast.Message{Nodes: []ast.Node{
			&ast.Text{Value: "'-'{-'-}'-'"},
		}}},
	} {
		tc := tc
		name := strconv.Itoa(idx)
		t.Run(name, func(t *testing.T) {
			require := require.New(t)

			r := strings.NewReader(tc.input)
			msg, err := Parse(r)
			require.NoError(err)
			require.Equal(tc.expected, msg)
		})
	}
}

func TestParseMessageText(t *testing.T) {
	for idx, tc := range []struct {
		input    string
		expected *ast.Text
	}{
		{"Spoon!",
			&ast.Text{Value: "Spoon!"}},
		{"Hello, {audience}!",
			&ast.Text{Value: "Hello, "}},
		{"It's peanut butter jelly time!",
			&ast.Text{Value: "It's peanut butter jelly time!"}},
		{"It''s peanut butter jelly time!",
			&ast.Text{Value: "It's peanut butter jelly time!"}},
		{"trailing quote'",
			&ast.Text{Value: "trailing quote'"}},
		{"-'{foo}-",
			&ast.Text{Value: "-{foo}-"}},
		{"-'{foo}'-",
			&ast.Text{Value: "-{foo}-"}},
		{"-'{foo}''-",
			&ast.Text{Value: "-{foo}'-"}},
		{"-'{foo}''-'",
			&ast.Text{Value: "-{foo}'-"}},
		{"-''{foo}''-",
			&ast.Text{Value: "-'"}},
		{"-'''{foo}'''-",
			&ast.Text{Value: "-'{foo}'-"}},
		{"'-{foo}-'",
			&ast.Text{Value: "'-"}},
		{"'-'{foo}'-'",
			&ast.Text{Value: "'-{foo}-'"}},
		{"'-''{foo}''-'",
			&ast.Text{Value: "'-'"}},
		{"'-'''{foo}'''-'",
			&ast.Text{Value: "'-'{foo}'-'"}},
	} {
		tc := tc
		label := strconv.Itoa(idx)
		t.Run(label, func(t *testing.T) {
			require := require.New(t)

			p := &parser{dec: NewDecoder(tc.input)}

			actual, err := p.parseMessageText()
			require.NoError(err)
			require.Equal(tc.expected, actual)
		})
	}
}
