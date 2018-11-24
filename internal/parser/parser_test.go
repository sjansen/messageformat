package parser

import (
	"strconv"
	"testing"

	"github.com/sjansen/messageformat/ast"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	for idx, tc := range []struct {
		pattern  string
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
		{"From: {begin}\nUntil: {end}", &ast.Message{Nodes: []ast.Node{
			&ast.Text{Value: "From: "},
			&ast.PlainArg{ArgID: "begin"},
			&ast.Text{Value: "\nUntil: "},
			&ast.PlainArg{ArgID: "end"},
		}}},
		{"From: {begin,date}\nUntil: {end,date,short}", &ast.Message{Nodes: []ast.Node{
			&ast.Text{Value: "From: "},
			&ast.SimpleArg{ArgID: "begin", ArgType: ast.DateType},
			&ast.Text{Value: "\nUntil: "},
			&ast.SimpleArg{ArgID: "end", ArgType: ast.DateType, ArgStyle: ast.ShortStyle},
		}}},
		{"{timespan, select, afternoon{Boa tarde, {name}.} evening{Boa noite, {name}.} other{Bom dia, {name}.}}",
			&ast.Message{Nodes: []ast.Node{
				&ast.SelectArg{ArgID: "timespan",
					Messages: map[string]*ast.Message{
						"afternoon": &ast.Message{Nodes: []ast.Node{
							&ast.Text{Value: "Boa tarde, "},
							&ast.PlainArg{ArgID: "name"},
							&ast.Text{Value: "."},
						}},
						"evening": &ast.Message{Nodes: []ast.Node{
							&ast.Text{Value: "Boa noite, "},
							&ast.PlainArg{ArgID: "name"},
							&ast.Text{Value: "."},
						}},
						"other": &ast.Message{Nodes: []ast.Node{
							&ast.Text{Value: "Bom dia, "},
							&ast.PlainArg{ArgID: "name"},
							&ast.Text{Value: "."},
						}},
					}},
			}}},
		{"{ timespan,select, afternoon{Boa tarde} evening{Boa noite} other{Bom dia} }, {name}.",
			&ast.Message{Nodes: []ast.Node{
				&ast.SelectArg{ArgID: "timespan",
					Messages: map[string]*ast.Message{
						"afternoon": &ast.Message{Nodes: []ast.Node{
							&ast.Text{Value: "Boa tarde"},
						}},
						"evening": &ast.Message{Nodes: []ast.Node{
							&ast.Text{Value: "Boa noite"},
						}},
						"other": &ast.Message{Nodes: []ast.Node{
							&ast.Text{Value: "Bom dia"},
						}},
					}},
				&ast.Text{Value: ", "},
				&ast.PlainArg{ArgID: "name"},
				&ast.Text{Value: "."},
			}}},
	} {
		tc := tc
		name := strconv.Itoa(idx)
		t.Run(name, func(t *testing.T) {
			require := require.New(t)

			msg, err := Parse(tc.pattern)
			require.NoError(err)
			require.Equal(tc.expected, msg)
		})
	}
}

func TestParseArgument(t *testing.T) {
	for idx, tc := range []struct {
		pattern  string
		expected ast.Node
	}{
		{"{foo}",
			&ast.PlainArg{ArgID: "foo"}},
		{"{ bar }",
			&ast.PlainArg{ArgID: "bar"}},
		{"{ baz } qux",
			&ast.PlainArg{ArgID: "baz"}},
		{"{0}",
			&ast.PlainArg{ArgID: "0"}},
		{"{1,duration}",
			&ast.SimpleArg{ArgID: "1", ArgType: ast.DurationType}},
		{"{ 2, number, percent }", &ast.SimpleArg{
			ArgID:    "2",
			ArgType:  ast.NumberType,
			ArgStyle: ast.PercentStyle}},
		{"{3,select,afternoon{Boa tarde!}evening{Boa noite!}other{Bom dia!}}", &ast.SelectArg{
			ArgID: "3",
			Messages: map[string]*ast.Message{
				"afternoon": &ast.Message{Nodes: []ast.Node{
					&ast.Text{Value: "Boa tarde!"},
				}},
				"evening": &ast.Message{Nodes: []ast.Node{
					&ast.Text{Value: "Boa noite!"},
				}},
				"other": &ast.Message{Nodes: []ast.Node{
					&ast.Text{Value: "Bom dia!"},
				}},
			}}},
	} {
		tc := tc
		label := strconv.Itoa(idx)
		t.Run(label, func(t *testing.T) {
			require := require.New(t)

			p := &parser{dec: NewDecoder(tc.pattern)}

			actual, err := p.parseArgument(0)
			require.NoError(err)
			require.Equal(tc.expected, actual)
		})
	}
}

func TestParseMessage(t *testing.T) {
	for idx, tc := range []struct {
		depth    int
		pattern  string
		expected []ast.Node
	}{
		{0, "Spoon!", []ast.Node{
			&ast.Text{Value: "Spoon!"},
		}},
		{1, "{Spoon!}", []ast.Node{
			&ast.Text{Value: "Spoon!"},
		}},
		{0, "Hello, {audience}!", []ast.Node{
			&ast.Text{Value: "Hello, "},
			&ast.PlainArg{ArgID: "audience"},
			&ast.Text{Value: "!"},
		}},
		{0, "{ greeting }, World!", []ast.Node{
			&ast.PlainArg{ArgID: "greeting"},
			&ast.Text{Value: ", World!"},
		}},
		{1, "{Hello, {audience}!}", []ast.Node{
			&ast.Text{Value: "Hello, "},
			&ast.PlainArg{ArgID: "audience"},
			&ast.Text{Value: "!"},
		}},
		{1, "{{ greeting }, World!}", []ast.Node{
			&ast.PlainArg{ArgID: "greeting"},
			&ast.Text{Value: ", World!"},
		}},
	} {
		tc := tc
		label := strconv.Itoa(idx)
		t.Run(label, func(t *testing.T) {
			require := require.New(t)

			p := &parser{dec: NewDecoder(tc.pattern)}

			actual, err := p.parseMessage(tc.depth)
			require.NoError(err)
			require.Equal(tc.expected, actual)
		})
	}
}

func TestParseMessageText(t *testing.T) {
	for idx, tc := range []struct {
		pattern  string
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

			p := &parser{dec: NewDecoder(tc.pattern)}

			actual, err := p.parseMessageText(0)
			require.NoError(err)
			require.Equal(tc.expected, actual)
		})
	}
}
