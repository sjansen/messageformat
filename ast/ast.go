package ast

//go:generate go run ../scripts/tmpl-to-go/main.go parts.go.tmpl

type Message struct {
	Parts []Part
}

type Part interface {
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
