package ast

type ArgType int

type ArgStyle int

type SimpleArg struct {
	Positions *Positions
	ArgID     string
	ArgType   ArgType
	ArgStyle  ArgStyle
}

const (
	DefaultType ArgType = iota
	DateType
	DurationType
	NumberType
	OrdinalType
	SpelloutType
	TimeType
	InvalidType
)

var argTypeFromKeyword = map[string]ArgType{
	"date":     DateType,
	"duration": DurationType,
	"number":   NumberType,
	"ordinal":  OrdinalType,
	"spellout": SpelloutType,
	"time":     TimeType,
}

func ArgTypeFromKeyword(keyword string) ArgType {
	if x, ok := argTypeFromKeyword[keyword]; ok {
		return x
	}
	return InvalidType
}

func (x ArgType) ToKeyword() string {
	switch x {
	case DateType:
		return "date"
	case DurationType:
		return "duration"
	case NumberType:
		return "number"
	case OrdinalType:
		return "ordinal"
	case SpelloutType:
		return "spellout"
	case TimeType:
		return "time"
	default:
		return ""
	}
}

const (
	DefaultStyle = iota
	CurrencyStyle
	FullStyle
	IntegerStyle
	LongStyle
	MediumStyle
	PercentStyle
	ShortStyle
	// non-keyword
	TextStyle
	SkeletonStyle
	InvalidStyle
)

var argStyleFromKeyword = map[string]ArgStyle{
	"currency": CurrencyStyle,
	"full":     FullStyle,
	"integer":  IntegerStyle,
	"long":     LongStyle,
	"medium":   MediumStyle,
	"percent":  PercentStyle,
	"short":    ShortStyle,
}

func ArgStyleFromKeyword(keyword string) ArgStyle {
	if x, ok := argStyleFromKeyword[keyword]; ok {
		return x
	}
	return InvalidStyle
}

func (x ArgStyle) ToKeyword() string {
	switch x {
	case CurrencyStyle:
		return "currency"
	case FullStyle:
		return "full"
	case IntegerStyle:
		return "integer"
	case LongStyle:
		return "long"
	case MediumStyle:
		return "medium"
	case PercentStyle:
		return "percent"
	case ShortStyle:
		return "short"
	default:
		return ""
	}
}
