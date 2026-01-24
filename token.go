package main

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL_TOK = "ILLEGAL"
	EOF_TOK     = "EOF"

	IDENT_TOK = "IDENT"
	INT_TOK   = "INT"

	ASSIGN_TOK    = "="
	PLUS_TOK      = "+"
	MINUS_TOK     = "-"
	ASTERISK_TOK  = "*"
	SLASH_TOK     = "/"
	BANG_TOK      = "!"
	EQUAL_TOK     = "=="
	NOT_EQUAL_TOK = "!="
	LT_TOK        = "<"
	GT_TOK        = ">"
	LTE_TOK       = "<="
	GTE_TOK       = ">="

	COMMA_TOK     = ","
	SEMICOLON_TOK = ";"

	LPAREN_TOK = "("
	RPAREN_TOK = ")"
	LBRACE_TOK = "{"
	RBRACE_TOK = "}"

	FN_TOK     = "FN"
	LET_TOK    = "LET"
	IF_TOK     = "IF"
	ELSE_TOK   = "ELSE"
	RETURN_TOK = "RETURN"
	TRUE_TOK   = "TRUE"
	FALSE_TOK  = "FALSE"
	NULL_TOK   = "NULL"
)

var keywords = map[string]TokenType{
	"let":    LET_TOK,
	"fn":     FN_TOK,
	"if":     IF_TOK,
	"else":   ELSE_TOK,
	"return": RETURN_TOK,
	"true":   TRUE_TOK,
	"false":  FALSE_TOK,
	"null":   NULL_TOK,
}

func NewToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT_TOK
}
