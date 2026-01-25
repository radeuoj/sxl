package main

import (
	"bytes"
	"strconv"
	"strings"
)

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type LetStatement struct {
	Name  *Identifier
	Value Expression
}

func (s *LetStatement) statementNode() {}

func (s *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString("let ")
	out.WriteString(s.Name.String())
	out.WriteString(" = ")

	if s.Value != nil {
		out.WriteString(s.Value.String())
	}

	out.WriteString(";\n")

	return out.String()
}

type ReturnStatement struct {
	Value Expression
}

func (s *ReturnStatement) statementNode() {}

func (s *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString("return ")
	out.WriteString(s.Value.String())
	out.WriteString(";\n")

	return out.String()
}

type BlockStatement struct {
	Statements []Statement
}

func (s *BlockStatement) statementNode() {}

func (s *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{\n")
	for _, s := range s.Statements {
		out.WriteString(s.String())
	}
	out.WriteString("}\n")

	return out.String()
}

type IfStatement struct {
	Condition Expression
	Then      Statement
	Else      Statement
}

func (s *IfStatement) statementNode() {}

func (s *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(s.Condition.String())
	out.WriteString(" ")
	out.WriteString(s.Then.String())

	if s.Else != nil {
		out.WriteString("else ")
		out.WriteString(s.Else.String())
	}

	return out.String()
}

type ExprStatement struct {
	Value Expression
}

func (s *ExprStatement) statementNode() {}

func (s *ExprStatement) String() string {
	return s.Value.String() + ";\n"
}

type FnStatement struct {
	Name  *Identifier
	Value *FnLiteral
}

func (s *FnStatement) statementNode() {}

func (s *FnStatement) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range s.Value.Params {
		params = append(params, p.String())
	}

	out.WriteString("fn ")
	out.WriteString(s.Name.Value)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(s.Value.Body.String())

	return out.String()
}

type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) String() string {
	return i.Value
}

type IntLiteral struct {
	Value int32
}

func (i *IntLiteral) expressionNode() {}

func (i *IntLiteral) String() string {
	return strconv.Itoa(int(i.Value))
}

type BoolLiteral struct {
	Value bool
}

func (b *BoolLiteral) expressionNode() {}

func (b *BoolLiteral) String() string {
	return strconv.FormatBool(b.Value)
}

type NullLiteral struct{}

func (b *NullLiteral) expressionNode() {}

func (b *NullLiteral) String() string {
	return "null"
}

type PrefixExpression struct {
	Operator string
	Right    Expression
}

func (e *PrefixExpression) expressionNode() {}

func (e *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(e.Operator)
	out.WriteString(e.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (e *InfixExpression) expressionNode() {}

func (e *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(e.Left.String())
	out.WriteString(" " + e.Operator + " ")
	out.WriteString(e.Right.String())
	out.WriteString(")")

	return out.String()
}

type CallExpression struct {
	Function  Expression
	Arguments []Expression
}

func (e *CallExpression) expressionNode() {}

func (e *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, arg := range e.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(e.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type FnLiteral struct {
	Params []*Identifier
	Body   *BlockStatement
}

func (e *FnLiteral) expressionNode() {}

func (e *FnLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range e.Params {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(e.Body.String())

	return out.String()
}
