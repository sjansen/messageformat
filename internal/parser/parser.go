package parser

import (
	"strings"
	"unicode/utf8"

	"github.com/sjansen/messageformat/ast"
	"github.com/sjansen/messageformat/errors"
)

func Parse(s string) (*ast.Message, error) {
	dec := NewDecoder(s)
	if nodes, err := parseMessage(dec, 0, false); err != nil {
		return nil, err
	} else {
		msg := &ast.Message{Nodes: nodes}
		return msg, nil
	}
}

func parseArgument(dec *Decoder, depth int) (ast.Node, error) {
	if err := requireRune(dec, '{'); err != nil {
		return nil, err
	}

	skipWhiteSpace(dec)
	argNameOrNumber := parseID(dec)
	skipWhiteSpace(dec)

	dec.Decode()
	ch := dec.Decoded()
	if ch == '}' {
		arg := &ast.PlainArg{ArgID: argNameOrNumber}
		return arg, nil
	} else if ch == ',' {
		skipWhiteSpace(dec)
	} else {
		return nil, &errors.UnexpectedToken{Token: string(ch)}
	}

	var arg ast.Node
	if keyword := parseID(dec); keyword == "select" {
		messages, err := parseSelectStyle(dec, depth)
		if err != nil {
			return nil, err
		}
		arg = &ast.SelectArg{ArgID: argNameOrNumber, Messages: messages}
	} else if keyword == "plural" || keyword == "selectordinal" {
		messages, err := parsePluralStyle(dec, depth)
		if err != nil {
			return nil, err
		}
		arg = &ast.PluralArg{
			ArgID:   argNameOrNumber,
			Ordinal: keyword == "selectordinal",
			// TODO Offset
			Messages: messages,
		}
	} else if argType := ast.ArgTypeFromKeyword(keyword); argType != ast.InvalidType {
		// TODO argStyleText and argSkeletonText
		argStyle, err := parseSimpleStyle(dec, depth)
		if err != nil {
			return nil, err
		}
		arg = &ast.SimpleArg{ArgID: argNameOrNumber, ArgType: argType, ArgStyle: argStyle}
	} else {
		return nil, &errors.UnexpectedToken{Token: keyword}
	}

	if err := requireRune(dec, '}'); err != nil {
		return nil, err
	}

	return arg, nil
}

func parseID(dec *Decoder) string {
	var b strings.Builder
	for dec.Decode() {
		ch := dec.Decoded()
		b.WriteRune(ch)
		next := dec.Peek()
		if isPatternWhiteSpace(next) || isPatternSyntax(next) {
			break
		}
	}
	return b.String()
}

func parseMessage(dec *Decoder, depth int, inPlural bool) ([]ast.Node, error) {
	nodes := []ast.Node{}
	if depth > 0 {
		if err := requireRune(dec, '{'); err != nil {
			return nil, err
		}
	}
	for {
		next := dec.Peek()
		if next == utf8.RuneError {
			break // TODO
		} else if depth > 0 && next == '}' {
			break
		} else if next == '{' {
			node, err := parseArgument(dec, depth)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		} else if inPlural && next == '#' {
			dec.Decode()
			nodes = append(nodes, &ast.PluralValue{})
		} else {
			node, err := parseMessageText(dec, depth, inPlural)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		}
	}
	if depth > 0 {
		if err := requireRune(dec, '}'); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func parseMessageText(dec *Decoder, depth int, inPlural bool) (*ast.Text, error) {
	inQuote := false
	var b strings.Builder
	for dec.Decode() {
		ch := dec.Decoded()
		if ch == '\'' {
			next := dec.Peek()
			if next == utf8.RuneError {
				if !inQuote {
					b.WriteRune('\'')
				}
				break
			} else if next == '\'' {
				b.WriteRune('\'')
				dec.Decode()
				if !inQuote {
					next := dec.Peek()
					if next == '{' || (depth > 0 && next == '}') {
						break
					} else if inPlural && next == '#' {
						break
					}
				}
			} else if inQuote {
				inQuote = false
			} else if next == '{' || next == '}' {
				inQuote = true
			} else if inPlural && next == '#' {
				inQuote = true
			} else {
				b.WriteRune('\'')
			}
		} else {
			b.WriteRune(ch)
			if !inQuote {
				next := dec.Peek()
				if next == '{' || (depth > 0 && next == '}') {
					break
				} else if inPlural && next == '#' {
					break
				}
			}
		}
	}
	t := &ast.Text{Value: b.String()}
	return t, nil
}

func parsePluralStyle(dec *Decoder, depth int) (map[string]*ast.Message, error) {
	skipWhiteSpace(dec)
	if err := requireRune(dec, ','); err != nil {
		return nil, err
	}

	messages := map[string]*ast.Message{}
	for {
		skipWhiteSpace(dec)
		next := dec.Peek()
		if next == '}' {
			return messages, nil
		}
		var id string
		if next == '=' {
			var b strings.Builder
			for dec.Decode() {
				ch := dec.Decoded()
				b.WriteRune(ch)
				if next := dec.Peek(); !isDigit(next) {
					break
				}
			}
			id = b.String()
		} else {
			id = parseID(dec)
		}
		skipWhiteSpace(dec)

		if nodes, err := parseMessage(dec, depth+1, true); err != nil {
			return nil, err
		} else {
			msg := &ast.Message{Nodes: nodes}
			messages[id] = msg
		}
	}
}

func parseSelectStyle(dec *Decoder, depth int) (map[string]*ast.Message, error) {
	skipWhiteSpace(dec)
	if err := requireRune(dec, ','); err != nil {
		return nil, err
	}

	messages := map[string]*ast.Message{}
	for {
		skipWhiteSpace(dec)
		next := dec.Peek()
		if next == '}' {
			return messages, nil
		}
		id := parseID(dec)
		skipWhiteSpace(dec)

		if nodes, err := parseMessage(dec, depth+1, false); err != nil {
			return nil, err
		} else {
			msg := &ast.Message{Nodes: nodes}
			messages[id] = msg
		}
	}
}

func parseSimpleStyle(dec *Decoder, depth int) (ast.ArgStyle, error) {
	skipWhiteSpace(dec)
	next := dec.Peek()
	if next == '}' {
		return ast.DefaultStyle, nil
	} else if next == ',' {
		dec.Decode()
	} else {
		return ast.DefaultStyle, &errors.UnexpectedToken{Token: string(next)}
	}

	skipWhiteSpace(dec)
	keyword := parseID(dec)
	argStyle := ast.ArgStyleFromKeyword(keyword)
	if argStyle == ast.InvalidStyle {
		return 0, &errors.UnexpectedToken{Token: keyword}
	}

	skipWhiteSpace(dec)
	return argStyle, nil
}

func requireRune(dec *Decoder, token rune) error {
	dec.Decode()
	ch := dec.Decoded()
	if ch == token {
		return nil
	}
	return &errors.UnexpectedToken{Token: string(ch)}
}

func skipWhiteSpace(dec *Decoder) {
	for next := dec.Peek(); isPatternWhiteSpace(next); next = dec.Peek() {
		if !dec.Decode() {
			break
		}
	}
}
