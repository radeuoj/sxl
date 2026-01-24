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

type ExprStatement struct {
	Value Expression
}

func (s *ExprStatement) statementNode() {}

func (s *ExprStatement) String() string {
	return s.Value.String() + ";\n"
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
