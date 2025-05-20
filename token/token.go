package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN = "ASSIGN"
	PLUS   = "PLUS"

	// Delimiters
	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"

	LPAREN = "LPAREN"
	RPAREN = "RPAREN"
	LBRACE = "LBRACE"
	RBRACE = "RBRACE"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

var keywords = map[string]TokenType{
	"let": LET,
	"fn":  FUNCTION,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
