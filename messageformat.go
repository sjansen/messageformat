package messageformat

import (
	"github.com/sjansen/messageformat/ast"
	"github.com/sjansen/messageformat/internal/parser"
)

func Parse(s string) (*ast.Message, error) {
	return parser.Parse(s)
}
