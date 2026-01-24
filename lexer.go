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
			tok.Type = EQUAL
			tok.Literal = "=="
		} else {
			tok = NewToken(ASSIGN, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = LTE
			tok.Literal = "<="
		} else {
			tok = NewToken(LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = GTE
			tok.Literal = ">="
		} else {
			tok = NewToken(GT, l.ch)
		}
	case '+':
		tok = NewToken(PLUS, l.ch)
	case '-':
		tok = NewToken(MINUS, l.ch)
	case '*':
		tok = NewToken(ASTERISK, l.ch)
	case '/':
		tok = NewToken(SLASH, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = EQUAL
			tok.Literal = "!="
		} else {
			tok = NewToken(BANG, l.ch)
		}
	case ',':
		tok = NewToken(COMMA, l.ch)
	case ';':
		tok = NewToken(SEMICOLON, l.ch)
	case '(':
		tok = NewToken(LPAREN, l.ch)
	case ')':
		tok = NewToken(RPAREN, l.ch)
	case '{':
		tok = NewToken(LBRACE, l.ch)
	case '}':
		tok = NewToken(RBRACE, l.ch)
	case 0:
		tok = NewToken(EOF, 0)
	default:
		if isIdentChar(l.ch) && !isDigit(l.ch) {
			tok.Literal = l.readIdent()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readInt()
			return tok
		} else {
			tok = NewToken(ILLEGAL, l.ch)
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
