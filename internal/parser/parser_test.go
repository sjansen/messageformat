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
		{"{0}",
			&ast.PlainArg{ArgID: "0"}},
		{"{foo}",
			&ast.PlainArg{ArgID: "foo"}},
		{"{ bar }",
			&ast.PlainArg{ArgID: "bar"}},
		{"{ baz } qux",
			&ast.PlainArg{ArgID: "baz"}},
		{"{4,duration}",
			&ast.SimpleArg{ArgID: "4", ArgType: ast.DurationType}},
		{"{ 5, number, percent }", &ast.SimpleArg{
			ArgID:    "5",
			ArgType:  ast.NumberType,
			ArgStyle: ast.PercentStyle}},
		{"{6,select,afternoon{Boa tarde!}evening{Boa noite!}other{Bom dia!}}", &ast.SelectArg{
			ArgID: "6",
			Messages: map[string]*ast.Message{
				"afternoon": &ast.Message{Nodes: []ast.Node{&ast.Text{Value: "Boa tarde!"}}},
				"evening":   &ast.Message{Nodes: []ast.Node{&ast.Text{Value: "Boa noite!"}}},
				"other":     &ast.Message{Nodes: []ast.Node{&ast.Text{Value: "Bom dia!"}}},
			}}},
		{"{7,plural,=0{no elves}one{one elf}other{multiple elves}}", &ast.PluralArg{
			ArgID: "7",
			Messages: map[string]*ast.Message{
				"=0":    &ast.Message{Nodes: []ast.Node{&ast.Text{Value: "no elves"}}},
				"one":   &ast.Message{Nodes: []ast.Node{&ast.Text{Value: "one elf"}}},
				"other": &ast.Message{Nodes: []ast.Node{&ast.Text{Value: "multiple elves"}}},
			}}},
		{`{ count, selectordinal,
		    one {#st item}
		    two {#nd item}
		    few {#rd item}
		    other {#th item} }`, &ast.PluralArg{
			ArgID:   "count",
			Ordinal: true,
			Messages: map[string]*ast.Message{
				"one": &ast.Message{Nodes: []ast.Node{
					&ast.PluralValue{},
					&ast.Text{Value: "st item"},
				}},
				"two": &ast.Message{Nodes: []ast.Node{
					&ast.PluralValue{},
					&ast.Text{Value: "nd item"},
				}},
				"few": &ast.Message{Nodes: []ast.Node{
					&ast.PluralValue{},
					&ast.Text{Value: "rd item"},
				}},
				"other": &ast.Message{Nodes: []ast.Node{
					&ast.PluralValue{},
					&ast.Text{Value: "th item"},
				}},
			}}},
	} {
		tc := tc
		label := strconv.Itoa(idx)
		t.Run(label, func(t *testing.T) {
			require := require.New(t)

			dec := NewDecoder(tc.pattern)

			actual, err := parseArgument(dec, 0)
			require.NoError(err)
			require.Equal(tc.expected, actual)
		})
	}
}

func TestParseMessage(t *testing.T) {
	for idx, tc := range []struct {
		depth    int
		inPlural bool
		pattern  string
		expected []ast.Node
	}{
		{0, false, "Spoon!", []ast.Node{
			&ast.Text{Value: "Spoon!"},
		}},
		{1, false, "{Spoon!}", []ast.Node{
			&ast.Text{Value: "Spoon!"},
		}},
		{0, false, "Hello, {audience}!", []ast.Node{
			&ast.Text{Value: "Hello, "},
			&ast.PlainArg{ArgID: "audience"},
			&ast.Text{Value: "!"},
		}},
		{0, false, "{ greeting }, World!", []ast.Node{
			&ast.PlainArg{ArgID: "greeting"},
			&ast.Text{Value: ", World!"},
		}},
		{1, false, "{Hello, {audience}!}", []ast.Node{
			&ast.Text{Value: "Hello, "},
			&ast.PlainArg{ArgID: "audience"},
			&ast.Text{Value: "!"},
		}},
		{1, false, "{{ greeting }, World!}", []ast.Node{
			&ast.PlainArg{ArgID: "greeting"},
			&ast.Text{Value: ", World!"},
		}},
		{1, false, "{The Internet is for #cats.}", []ast.Node{
			&ast.Text{Value: "The Internet is for #cats."},
		}},
		{1, true, "{# {color} items}", []ast.Node{
			&ast.PluralValue{},
			&ast.Text{Value: " "},
			&ast.PlainArg{ArgID: "color"},
			&ast.Text{Value: " items"},
		}},
	} {
		tc := tc
		label := strconv.Itoa(idx)
		t.Run(label, func(t *testing.T) {
			require := require.New(t)

			dec := NewDecoder(tc.pattern)

			actual, err := parseMessage(dec, tc.depth, tc.inPlural)
			require.NoError(err)
			require.Equal(tc.expected, actual)
		})
	}
}

func TestParseMessageText(t *testing.T) {
	for idx, tc := range []struct {
		inPlural bool
		pattern  string
		expected *ast.Text
	}{
		{false, "Spoon!",
			&ast.Text{Value: "Spoon!"}},
		{false, "Hello, {audience}!",
			&ast.Text{Value: "Hello, "}},
		{false, "It's peanut butter jelly time!",
			&ast.Text{Value: "It's peanut butter jelly time!"}},
		{false, "It''s peanut butter jelly time!",
			&ast.Text{Value: "It's peanut butter jelly time!"}},
		{false, "trailing quote'",
			&ast.Text{Value: "trailing quote'"}},
		{false, "-'{foo}-",
			&ast.Text{Value: "-{foo}-"}},
		{false, "-'{foo}'-",
			&ast.Text{Value: "-{foo}-"}},
		{false, "-'{foo}''-",
			&ast.Text{Value: "-{foo}'-"}},
		{false, "-'{foo}''-'",
			&ast.Text{Value: "-{foo}'-"}},
		{false, "-''{foo}''-",
			&ast.Text{Value: "-'"}},
		{false, "-'''{foo}'''-",
			&ast.Text{Value: "-'{foo}'-"}},
		{false, "'-{foo}-'",
			&ast.Text{Value: "'-"}},
		{false, "'-'{foo}'-'",
			&ast.Text{Value: "'-{foo}-'"}},
		{false, "'-''{foo}''-'",
			&ast.Text{Value: "'-'"}},
		{false, "'-'''{foo}'''-'",
			&ast.Text{Value: "'-'{foo}'-'"}},
		{false, "We're #1!",
			&ast.Text{Value: "We're #1!"}},
		{true, "count: #",
			&ast.Text{Value: "count: "}},
		{true, "-'#'-",
			&ast.Text{Value: "-#-"}},
		{false, "-'#'-",
			&ast.Text{Value: "-'#'-"}},
		{false, "'{{ foo }}'",
			&ast.Text{Value: "{{ foo }}"}},
		{true, "'{# foo #}'",
			&ast.Text{Value: "{# foo #}"}},
	} {
		tc := tc
		label := strconv.Itoa(idx)
		t.Run(label, func(t *testing.T) {
			require := require.New(t)

			dec := NewDecoder(tc.pattern)

			actual, err := parseMessageText(dec, 0, tc.inPlural)
			require.NoError(err)
			require.Equal(tc.expected, actual)
		})
	}
}
