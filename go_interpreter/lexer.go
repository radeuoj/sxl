package main

type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPos]
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = EQUAL_TOK
			tok.Literal = "=="
		} else {
			tok = NewToken(ASSIGN_TOK, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = LTE_TOK
			tok.Literal = "<="
		} else {
			tok = NewToken(LT_TOK, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = GTE_TOK
			tok.Literal = ">="
		} else {
			tok = NewToken(GT_TOK, l.ch)
		}
	case '+':
		tok = NewToken(PLUS_TOK, l.ch)
	case '-':
		tok = NewToken(MINUS_TOK, l.ch)
	case '*':
		tok = NewToken(ASTERISK_TOK, l.ch)
	case '/':
		if l.peekChar() == '/' { // comment
			l.readChar()
			l.readChar()
			l.skipComment()
			return l.NextToken()
		} else {
			tok = NewToken(SLASH_TOK, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = NOT_EQUAL_TOK
			tok.Literal = "!="
		} else {
			tok = NewToken(BANG_TOK, l.ch)
		}
	case ',':
		tok = NewToken(COMMA_TOK, l.ch)
	case ';':
		tok = NewToken(SEMICOLON_TOK, l.ch)
	case '(':
		tok = NewToken(LPAREN_TOK, l.ch)
	case ')':
		tok = NewToken(RPAREN_TOK, l.ch)
	case '{':
		tok = NewToken(LBRACE_TOK, l.ch)
	case '}':
		tok = NewToken(RBRACE_TOK, l.ch)
	case '"':
		tok.Type = STRING_TOK
		tok.Literal = l.readString()

		if l.ch == 0 {
			return NewToken(ILLEGAL_TOK, '"')
		}
	case 0:
		tok = NewToken(EOF_TOK, 0)
	default:
		if isIdentChar(l.ch) && !isDigit(l.ch) {
			tok.Literal = l.readIdent()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = INT_TOK
			tok.Literal = l.readInt()
			return tok
		} else {
			tok = NewToken(ILLEGAL_TOK, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) readIdent() string {
	pos := l.pos

	for isIdentChar(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.pos]
}

func (l *Lexer) readInt() string {
	pos := l.pos

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.pos]
}

func isIdentChar(ch byte) bool {
	return ('a' <= ch && ch <= 'z') ||
		('A' <= ch && ch <= 'Z') ||
		isDigit(ch) || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readString() string {
	pos := l.pos + 1

	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[pos:l.pos]
}
