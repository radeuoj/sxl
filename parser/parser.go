package parser

import (
	"fmt"
	"mkl/ast"
	"mkl/lexer"
	"mkl/token"
)

type Parser struct {
	l         *lexer.Lexer
	errors    []string
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(tok token.TokenType) {
	msg := fmt.Sprintf("expected %s but got %s", tok, p.peekToken.Literal)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.isCurToken(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.isCurToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{}
	p.nextToken()

	for !p.isCurToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) isCurToken(tok token.TokenType) bool {
	return p.curToken.Type == tok
}

func (p *Parser) isPeekToken(tok token.TokenType) bool {
	return p.peekToken.Type == tok
}

func (p *Parser) expectPeek(tok token.TokenType) bool {
	if p.isPeekToken(tok) {
		p.nextToken()
		return true
	} else {
		p.peekError(tok)
		return false
	}
}
