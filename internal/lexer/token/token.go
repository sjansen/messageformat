package token

const (
	TEXT   TokenType = "text"
	COMMA  TokenType = ","
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"
	EOF    TokenType = "eof"
)

type TokenType string

type Token struct {
	Type  TokenType
	Value string
}
