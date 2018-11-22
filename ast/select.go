package ast

type SelectArg struct {
	Positions *Positions
	ArgID      string
	Nodes      map[string]Node
}
