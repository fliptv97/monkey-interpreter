package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	Illegal = "Illegal"
	EOF     = "EOF"

	Ident = "Ident"
	Int   = "Int"

	Assign   = "="
	Plus     = "+"
	Minus    = "-"
	Bang     = "!"
	Asterisk = "*"
	Slash    = "/"

	LessThan    = "<"
	GreaterThan = ">"

	Equal    = "=="
	NotEqual = "!="

	Comma     = ","
	Semicolon = ";"

	LeftParen  = "("
	RightParen = ")"
	LeftBrace  = "{"
	RightBrace = "}"

	Function = "Function"
	Let      = "Let"
	True     = "True"
	False    = "False"
	If       = "If"
	Else     = "Else"
	Return   = "Return"
)

var keywords = map[string]Type{
	"let":    Let,
	"fn":     Function,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}
