package parser

func isPatternSyntax(r rune) bool {
	switch {
	case r >= '\u0021' && r <= '\u002F':
		fallthrough
	case r >= '\u003A' && r <= '\u0040':
		fallthrough
	case r >= '\u005B' && r <= '\u005E':
		fallthrough
	case r == '\u0060':
		fallthrough
	case r >= '\u007B' && r <= '\u007E':
		fallthrough
	case r >= '\u00A1' && r <= '\u00A7':
		fallthrough
	case r == '\u00A9' || r == '\u00AB' || r == '\u00AC' || r == '\u00AE':
		fallthrough
	case r == '\u00B0' || r == '\u00B1' || r == '\u00B6' || r == '\u00BB' || r == '\u00BF':
		fallthrough
	case r == '\u00D7' || r == '\u00F7':
		fallthrough
	case r >= '\u2010' && r <= '\u2027':
		fallthrough
	case r >= '\u2030' && r <= '\u203E':
		fallthrough
	case r >= '\u2041' && r <= '\u2053':
		fallthrough
	case r >= '\u2055' && r <= '\u205E':
		fallthrough
	case r >= '\u2190' && r <= '\u245F':
		fallthrough
	case r >= '\u2500' && r <= '\u2775':
		fallthrough
	case r >= '\u2794' && r <= '\u2BFF':
		fallthrough
	case r >= '\u2E00' && r <= '\u2EF7':
		fallthrough
	case r >= '\u3001' && r <= '\u3003':
		fallthrough
	case r >= '\u3008' && r <= '\u3020':
		fallthrough
	case r == '\u3030':
		fallthrough
	case r == '\uFD3E' || r == '\uFD3F':
		fallthrough
	case r == '\uFE45' || r == '\uFE46':
		return true
	default:
		return false
	}
}

func isPatternWhiteSpace(r rune) bool {
	switch {
	case r == ' ': // U+0020
		fallthrough
	case r >= '\u0009' && r <= '\u000D':
		fallthrough
	case r == '\u0085':
		fallthrough
	case r == '\u200E' || r == '\u200F':
		fallthrough
	case r == '\u2028' || r == '\u2029':
		return true
	default:
		return false
	}
}
