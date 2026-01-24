package main

import (
	"fmt"
	"strconv"
)

type (
	PrefixParseFn func() Expression
	InfixParseFn  func(Expression) Expression
)

const (
	LOWEST_P int = iota
	ASSIGN_P
	EQUALS_P
	SUM_P
	PRODUCT_P
	PREFIX_P
	CALL_P
)

var precedences = map[TokenType]int{
	EQUAL_TOK:     EQUALS_P,
	NOT_EQUAL_TOK: EQUALS_P,
	LT_TOK:        EQUALS_P,
	LTE_TOK:       EQUALS_P,
	GT_TOK:        EQUALS_P,
	GTE_TOK:       EQUALS_P,
	PLUS_TOK:      SUM_P,
	MINUS_TOK:     SUM_P,
	ASTERISK_TOK:  PRODUCT_P,
	SLASH_TOK:     PRODUCT_P,
	ASSIGN_TOK:    ASSIGN_P,
	LPAREN_TOK:    CALL_P,
}

type Parser struct {
	l         *Lexer
	errors    []string
	curToken  Token
	peekToken Token

	prefixParseFns map[TokenType]PrefixParseFn
	infixParseFns  map[TokenType]InfixParseFn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[TokenType]PrefixParseFn)
	p.registerPrefix(IDENT_TOK, p.parseIdentifier)
	p.registerPrefix(INT_TOK, p.parseIntLiteral)
	p.registerPrefix(TRUE_TOK, p.parseBoolLiteral)
	p.registerPrefix(FALSE_TOK, p.parseBoolLiteral)
	p.registerPrefix(NULL_TOK, p.parseNullLiteral)
	p.registerPrefix(MINUS_TOK, p.parsePrefixExpression)
	p.registerPrefix(BANG_TOK, p.parsePrefixExpression)
	p.registerPrefix(LPAREN_TOK, p.parseGroupedExpression)

	p.infixParseFns = make(map[TokenType]InfixParseFn)
	p.registerInfix(EQUAL_TOK, p.parseInfixExpression)
	p.registerInfix(NOT_EQUAL_TOK, p.parseInfixExpression)
	p.registerInfix(LT_TOK, p.parseInfixExpression)
	p.registerInfix(LTE_TOK, p.parseInfixExpression)
	p.registerInfix(GT_TOK, p.parseInfixExpression)
	p.registerInfix(GTE_TOK, p.parseInfixExpression)
	p.registerInfix(PLUS_TOK, p.parseInfixExpression)
	p.registerInfix(MINUS_TOK, p.parseInfixExpression)
	p.registerInfix(ASTERISK_TOK, p.parseInfixExpression)
	p.registerInfix(SLASH_TOK, p.parseInfixExpression)
	p.registerInfix(ASSIGN_TOK, p.parseInfixExpression)
	p.registerInfix(LPAREN_TOK, p.parseCallExpression)

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(tok TokenType) {
	msg := fmt.Sprintf("expected %s but got %s", tok, p.peekToken.Literal)
	p.errors = append(p.errors, msg)
}

func (p *Parser) invalidPrefixOperatorError(tok TokenType) {
	msg := fmt.Sprintf("invalid prefix operator %s", tok)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) isCurToken(tok TokenType) bool {
	return p.curToken.Type == tok
}

func (p *Parser) isPeekToken(tok TokenType) bool {
	return p.peekToken.Type == tok
}

func (p *Parser) expectPeek(tok TokenType) bool {
	if p.isPeekToken(tok) {
		p.nextToken()
		return true
	} else {
		p.peekError(tok)
		return false
	}
}

func (p *Parser) registerPrefix(tok TokenType, fn PrefixParseFn) {
	p.prefixParseFns[tok] = fn
}

func (p *Parser) registerInfix(tok TokenType, fn InfixParseFn) {
	p.infixParseFns[tok] = fn
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST_P
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST_P
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for !p.isCurToken(EOF_TOK) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

// returns statement and if it ends with a semicolon
func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case LET_TOK:
		return p.parseLetStatement()
	case RETURN_TOK:
		return p.parseReturnStatement()
	case LBRACE_TOK:
		return p.parseBlockStatement()
	case IF_TOK:
		return p.parseIfStatement()
	default:
		return p.parseExprStatement()
	}
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{}

	if !p.expectPeek(IDENT_TOK) {
		return nil
	}

	stmt.Name = &Identifier{Value: p.curToken.Literal}

	if !p.expectPeek(ASSIGN_TOK) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST_P)

	if !p.expectPeek(SEMICOLON_TOK) {
		return nil
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST_P)

	if !p.expectPeek(SEMICOLON_TOK) {
		return nil
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Statements: []Statement{}}
	p.nextToken()

	for !p.isCurToken(RBRACE_TOK) && !p.isCurToken(EOF_TOK) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
		p.nextToken()
	}

	return block
}

func (p *Parser) parseIfStatement() *IfStatement {
	stmt := &IfStatement{}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST_P)

	if !p.expectPeek(LBRACE_TOK) {
		return nil
	}
	stmt.Then = p.parseBlockStatement()

	if p.isPeekToken(ELSE_TOK) {
		p.nextToken()

		if p.isPeekToken(IF_TOK) { // if else syntax
			p.nextToken()
			stmt.Else = p.parseIfStatement()
		} else if p.expectPeek(LBRACE_TOK) {
			stmt.Else = p.parseBlockStatement()
		} else {
			return nil
		}
	}

	return stmt
}

// returns statement and if it ends with a semicolon
func (p *Parser) parseExprStatement() *ExprStatement {
	stmt := &ExprStatement{}
	stmt.Value = p.parseExpression(LOWEST_P)

	if !p.expectPeek(SEMICOLON_TOK) {
		return nil
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.invalidPrefixOperatorError(p.curToken.Type)
		return nil
	}
	left := prefix()

	for !p.isPeekToken(SEMICOLON_TOK) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return nil
		}

		p.nextToken()
		left = infix(left)
	}

	return left
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Value: p.curToken.Literal}
}

func (p *Parser) parseIntLiteral() Expression {
	lit := &IntLiteral{}

	val, err := strconv.ParseInt(p.curToken.Literal, 0, 32)
	if err != nil {
		msg := fmt.Sprintf("failed to parse %s to int", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = int32(val)
	return lit
}

func (p *Parser) parseBoolLiteral() Expression {
	return &BoolLiteral{Value: p.isCurToken(TRUE_TOK)}
}

func (p *Parser) parseNullLiteral() Expression {
	return &NullLiteral{}
}

func (p *Parser) parsePrefixExpression() Expression {
	expr := &PrefixExpression{
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expr.Right = p.parseExpression(PREFIX_P)

	return expr
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expr := &InfixExpression{
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()
	expr := p.parseExpression(LOWEST_P)

	if !p.expectPeek(RPAREN_TOK) {
		return nil
	}

	return expr
}

func (p *Parser) parseCallExpression(fn Expression) Expression {
	return &CallExpression{
		Function:  fn,
		Arguments: p.parseCallArguments(),
	}
}

func (p *Parser) parseCallArguments() []Expression {
	args := []Expression{}

	if p.isPeekToken(RPAREN_TOK) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST_P))

	for p.isPeekToken(COMMA_TOK) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST_P))
	}

	if !p.expectPeek(RPAREN_TOK) {
		return nil
	}

	return args
}
