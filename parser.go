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
	LOWEST = iota
	EQUALS
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[TokenType]int{
	EQUAL:     EQUALS,
	NOT_EQUAL: EQUALS,
	LT:        EQUALS,
	LTE:       EQUALS,
	GT:        EQUALS,
	GTE:       EQUALS,
	PLUS:      SUM,
	MINUS:     SUM,
	ASTERISK:  PRODUCT,
	SLASH:     PRODUCT,
	LPAREN:    CALL,
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
	p.registerPrefix(IDENT, p.parseIdentifier)
	p.registerPrefix(INT, p.parseIntLiteral)
	p.registerPrefix(TRUE, p.parseBoolLiteral)
	p.registerPrefix(FALSE, p.parseBoolLiteral)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(BANG, p.parsePrefixExpression)
	p.registerPrefix(LPAREN, p.parseGroupedExpression)

	p.infixParseFns = make(map[TokenType]InfixParseFn)
	p.registerInfix(EQUAL, p.parseInfixExpression)
	p.registerInfix(NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(LT, p.parseInfixExpression)
	p.registerInfix(LTE, p.parseInfixExpression)
	p.registerInfix(GT, p.parseInfixExpression)
	p.registerInfix(GTE, p.parseInfixExpression)
	p.registerInfix(PLUS, p.parseInfixExpression)
	p.registerInfix(MINUS, p.parseInfixExpression)
	p.registerInfix(ASTERISK, p.parseInfixExpression)
	p.registerInfix(SLASH, p.parseInfixExpression)
	p.registerInfix(LPAREN, p.parseCallExpression)

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

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for !p.isCurToken(EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case LET:
		return p.parseLetStatement()
	case RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExprStatement()
	}
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Value: p.curToken.Literal}

	if !p.expectPeek(ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.isPeekToken(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.isPeekToken(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExprStatement() *ExprStatement {
	stmt := &ExprStatement{}
	stmt.Value = p.parseExpression(LOWEST)

	if p.isPeekToken(SEMICOLON) {
		p.nextToken()
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

	for !p.isPeekToken(SEMICOLON) && precedence < p.peekPrecedence() {
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
	return &BoolLiteral{Value: p.isCurToken(TRUE)}
}

func (p *Parser) parsePrefixExpression() Expression {
	expr := &PrefixExpression{
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)

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
	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
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

	if p.isPeekToken(RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.isPeekToken(COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return args
}
