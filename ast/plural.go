package ast

type PluralArg struct {
	Positions *Positions
	ArgID     string
	Ordinal   bool
	Offset    int
	Nodes     map[string]Node
}
