// Code generated by scripts/tmpl-to-go. DO NOT EDIT.
package ast

/*
type Argument struct {
	ArgID  string
}
*/

var _ Argument = &PlainArg{}

func (x *PlainArg) ArgNameOrNumber() string {
	return x.ArgID
}

var _ Argument = &PluralArg{}

func (x *PluralArg) ArgNameOrNumber() string {
	return x.ArgID
}

var _ Argument = &SelectArg{}

func (x *SelectArg) ArgNameOrNumber() string {
	return x.ArgID
}

var _ Argument = &SimpleArg{}

func (x *SimpleArg) ArgNameOrNumber() string {
	return x.ArgID
}
