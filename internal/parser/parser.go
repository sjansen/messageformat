package parser

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/sjansen/messageformat/ast"
	"github.com/sjansen/messageformat/errors"
	"github.com/sjansen/messageformat/internal/decoder"
)

func Parse(s string) (*ast.Message, error) {
	dec := decoder.New(s)
	parts, err := parseMessage(dec, 0, false)
	if err != nil {
		return nil, err
	}
	msg := &ast.Message{Parts: parts}
	return msg, nil
}

func parseArgument(dec *decoder.Decoder, depth int) (ast.Part, error) {
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

	var arg ast.Part
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

func parseID(dec *decoder.Decoder) string {
	var b strings.Builder
	for dec.Decode() {
		ch := dec.Decoded()
		b.WriteRune(ch)
		next := dec.Peek()
		if unicode.In(next, unicode.Pattern_White_Space, unicode.Pattern_Syntax) {
			break
		}
	}
	return b.String()
}

func parseMessage(dec *decoder.Decoder, depth int, inPlural bool) ([]ast.Part, error) {
	parts := []ast.Part{}
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
			part, err := parseArgument(dec, depth)
			if err != nil {
				return nil, err
			}
			parts = append(parts, part)
		} else if inPlural && next == '#' {
			dec.Decode()
			parts = append(parts, &ast.NumberSign{})
		} else {
			part, err := parseMessageText(dec, depth, inPlural)
			if err != nil {
				return nil, err
			}
			parts = append(parts, part)
		}
	}
	if depth > 0 {
		if err := requireRune(dec, '}'); err != nil {
			return nil, err
		}
	}
	return parts, nil
}

func parseMessageText(dec *decoder.Decoder, depth int, inPlural bool) (*ast.Text, error) {
	var b strings.Builder
	for dec.Decode() {
		ch := dec.Decoded()
		if ch == '\'' {
			done := parseMessageTextAfterQuote(&b, dec, depth, inPlural)
			if done {
				break
			}
		} else {
			b.WriteRune(ch)
			next := dec.Peek()
			if next == '{' || (depth > 0 && next == '}') || (inPlural && next == '#') {
				break
			}
		}
	}
	t := &ast.Text{Value: b.String()}
	return t, nil
}

func parseMessageTextAfterQuote(b *strings.Builder, dec *decoder.Decoder, depth int, inPlural bool) bool {
	done := false
	next := dec.Peek()
	if next == utf8.RuneError {
		b.WriteRune('\'')
		done = true
	} else if next == '\'' {
		b.WriteRune('\'')
		dec.Decode()
		next := dec.Peek()
		if next == '{' || (depth > 0 && next == '}') || (inPlural && next == '#') {
			done = true
		}
	} else if next == '{' || next == '}' || (inPlural && next == '#') {
		parseMessageTextInQuote(b, dec)
	} else {
		b.WriteRune('\'')
	}
	return done
}

func parseMessageTextInQuote(b *strings.Builder, dec *decoder.Decoder) {
	for dec.Decode() {
		ch := dec.Decoded()
		if ch != '\'' {
			b.WriteRune(ch)
		} else {
			next := dec.Peek()
			if next == utf8.RuneError {
				break
			} else if next == '\'' {
				b.WriteRune('\'')
				dec.Decode()
			} else {
				return
			}
		}
	}
}

func parsePluralStyle(dec *decoder.Decoder, depth int) (map[string]*ast.Message, error) {
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
				if next := dec.Peek(); next < '0' || next > '9' {
					break
				}
			}
			id = b.String()
		} else {
			id = parseID(dec)
		}
		skipWhiteSpace(dec)

		parts, err := parseMessage(dec, depth+1, true)
		if err != nil {
			return nil, err
		}
		msg := &ast.Message{Parts: parts}
		messages[id] = msg
	}
}

func parseSelectStyle(dec *decoder.Decoder, depth int) (map[string]*ast.Message, error) {
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

		parts, err := parseMessage(dec, depth+1, false)
		if err != nil {
			return nil, err
		}
		msg := &ast.Message{Parts: parts}
		messages[id] = msg
	}
}

func parseSimpleStyle(dec *decoder.Decoder, depth int) (ast.ArgStyle, error) {
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

func requireRune(dec *decoder.Decoder, token rune) error {
	dec.Decode()
	ch := dec.Decoded()
	if ch == token {
		return nil
	}
	return &errors.UnexpectedToken{Token: string(ch)}
}

func skipWhiteSpace(dec *decoder.Decoder) {
	for next := dec.Peek(); unicode.In(next, unicode.Pattern_White_Space); next = dec.Peek() {
		if !dec.Decode() {
			break
		}
	}
}
