package token

const (
	TEXT   TokenType = "text"
	COMMA  TokenType = ","
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"
)

type TokenType string

type Token struct {
	Type  TokenType
	Value string
}
