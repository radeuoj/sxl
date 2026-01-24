package main

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	ASTERISK  = "*"
	SLASH     = "/"
	BANG      = "!"
	EQUAL     = "=="
	NOT_EQUAL = "!="
	LT        = "<"
	GT        = ">"
	LTE       = "<="
	GTE       = ">="

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FN     = "FN"
	LET    = "LET"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
)

var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FN,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func NewToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
