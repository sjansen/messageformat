package ast

//go:generate go run ../scripts/tmpl-to-go/main.go nodes.go.tmpl

type Message struct {
	Nodes []Node
}

type Node interface {
	HasPositions() bool
	Begin() Position
	End() Position
}

type Positions struct {
	Begin Position
	End   Position
}

type Position struct {
	Line       int
	ByteColumn int
	RuneColumn int
}
