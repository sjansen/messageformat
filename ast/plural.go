package ast

type PluralArg struct {
	Positions *Positions
	ArgID     string
	Ordinal   bool
	Offset    int
	Messages  map[string]*Message
}

type PluralValue struct {
	Positions *Positions
}
