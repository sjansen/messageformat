package ast

type SelectArg struct {
	Positions *Positions
	ArgID     string
	Messages  map[string]*Message
}
