package parser

import (
	"strings"
	"unicode/utf8"

	"github.com/sjansen/messageformat/ast"
	"github.com/sjansen/messageformat/errors"
)

func Parse(s string) (*ast.Message, error) {
	p := &parser{dec: NewDecoder(s)}
	if nodes, err := p.parseMessage(0, false); err != nil {
		return nil, err
	} else {
		msg := &ast.Message{Nodes: nodes}
		return msg, nil
	}
}

type parser struct {
	dec *Decoder
}

func (p *parser) parseArgument(depth int) (ast.Node, error) {
	if err := p.requireRune('{'); err != nil {
		return nil, err
	}

	p.skipWhiteSpace()
	argNameOrNumber := p.parseID()
	p.skipWhiteSpace()

	p.dec.Decode()
	ch := p.dec.Decoded()
	if ch == '}' {
		arg := &ast.PlainArg{ArgID: argNameOrNumber}
		return arg, nil
	} else if ch == ',' {
		p.skipWhiteSpace()
	} else {
		return nil, &errors.UnexpectedToken{Token: string(ch)}
	}

	var arg ast.Node
	if keyword := p.parseID(); keyword == "select" {
		messages, err := p.parseSelectStyle(depth)
		if err != nil {
			return nil, err
		}
		arg = &ast.SelectArg{ArgID: argNameOrNumber, Messages: messages}
	} else if keyword == "plural" || keyword == "selectordinal" {
		messages, err := p.parsePluralStyle(depth)
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
		argStyle, err := p.parseSimpleStyle(depth)
		if err != nil {
			return nil, err
		}
		arg = &ast.SimpleArg{ArgID: argNameOrNumber, ArgType: argType, ArgStyle: argStyle}
	} else {
		return nil, &errors.UnexpectedToken{Token: keyword}
	}

	if err := p.requireRune('}'); err != nil {
		return nil, err
	}

	return arg, nil
}

func (p *parser) parseID() string {
	var b strings.Builder
	for p.dec.Decode() {
		ch := p.dec.Decoded()
		b.WriteRune(ch)
		next := p.dec.Peek()
		if isPatternWhiteSpace(next) || isPatternSyntax(next) {
			break
		}
	}
	return b.String()
}

func (p *parser) parseMessage(depth int, inPlural bool) ([]ast.Node, error) {
	nodes := []ast.Node{}
	if depth > 0 {
		if err := p.requireRune('{'); err != nil {
			return nil, err
		}
	}
	for {
		next := p.dec.Peek()
		if next == utf8.RuneError {
			break // TODO
		} else if depth > 0 && next == '}' {
			break
		} else if next == '{' {
			node, err := p.parseArgument(depth)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		} else if inPlural && next == '#' {
			p.dec.Decode()
			nodes = append(nodes, &ast.PluralValue{})
		} else {
			node, err := p.parseMessageText(depth, inPlural)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		}
	}
	if depth > 0 {
		if err := p.requireRune('}'); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (p *parser) parseMessageText(depth int, inPlural bool) (*ast.Text, error) {
	inQuote := false
	var b strings.Builder
	for p.dec.Decode() {
		ch := p.dec.Decoded()
		if ch == '\'' {
			next := p.dec.Peek()
			if next == utf8.RuneError {
				if !inQuote {
					b.WriteRune('\'')
				}
				break
			} else if next == '\'' {
				b.WriteRune('\'')
				p.dec.Decode()
				if !inQuote {
					next := p.dec.Peek()
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
				next := p.dec.Peek()
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

func (p *parser) parsePluralStyle(depth int) (map[string]*ast.Message, error) {
	p.skipWhiteSpace()
	if err := p.requireRune(','); err != nil {
		return nil, err
	}

	messages := map[string]*ast.Message{}
	for {
		p.skipWhiteSpace()
		next := p.dec.Peek()
		if next == '}' {
			return messages, nil
		}
		var id string
		if next == '=' {
			var b strings.Builder
			for p.dec.Decode() {
				ch := p.dec.Decoded()
				b.WriteRune(ch)
				if next := p.dec.Peek(); !isDigit(next) {
					break
				}
			}
			id = b.String()
		} else {
			id = p.parseID()
		}
		p.skipWhiteSpace()

		if nodes, err := p.parseMessage(depth+1, true); err != nil {
			return nil, err
		} else {
			msg := &ast.Message{Nodes: nodes}
			messages[id] = msg
		}
	}
}

func (p *parser) parseSelectStyle(depth int) (map[string]*ast.Message, error) {
	p.skipWhiteSpace()
	if err := p.requireRune(','); err != nil {
		return nil, err
	}

	messages := map[string]*ast.Message{}
	for {
		p.skipWhiteSpace()
		next := p.dec.Peek()
		if next == '}' {
			return messages, nil
		}
		id := p.parseID()
		p.skipWhiteSpace()

		if nodes, err := p.parseMessage(depth+1, false); err != nil {
			return nil, err
		} else {
			msg := &ast.Message{Nodes: nodes}
			messages[id] = msg
		}
	}
}

func (p *parser) parseSimpleStyle(depth int) (ast.ArgStyle, error) {
	p.skipWhiteSpace()
	next := p.dec.Peek()
	if next == '}' {
		return ast.DefaultStyle, nil
	} else if next == ',' {
		p.dec.Decode()
	} else {
		return ast.DefaultStyle, &errors.UnexpectedToken{Token: string(next)}
	}

	p.skipWhiteSpace()
	keyword := p.parseID()
	argStyle := ast.ArgStyleFromKeyword(keyword)
	if argStyle == ast.InvalidStyle {
		return 0, &errors.UnexpectedToken{Token: keyword}
	}

	p.skipWhiteSpace()
	return argStyle, nil
}

func (p *parser) requireRune(token rune) error {
	p.dec.Decode()
	ch := p.dec.Decoded()
	if ch == token {
		return nil
	}
	return &errors.UnexpectedToken{Token: string(ch)}
}

func (p *parser) skipWhiteSpace() {
	for next := p.dec.Peek(); isPatternWhiteSpace(next); next = p.dec.Peek() {
		if !p.dec.Decode() {
			break
		}
	}
}
