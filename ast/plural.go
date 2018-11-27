package ast

type NumberSign struct {
	Positions *Positions
}

type PluralArg struct {
	Positions *Positions
	ArgID     string
	Ordinal   bool
	Offset    int
	Messages  map[string]*Message
}
