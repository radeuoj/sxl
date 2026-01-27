package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type ValueType string

const (
	INT_VAL_T      = "int"
	BOOL_VAL_T     = "bool"
	NULL_VAL_T     = "null"
	ERROR_VAL_T    = "error"
	FN_VAL_T       = "fn"
	RETURN_VAL_T   = "return"
	STRING_VAL_T   = "string"
	BUILTIN_FN_VAL = "builtin fn"
)

var (
	TRUE_VAL  = &BoolValue{Value: true}
	FALSE_VAL = &BoolValue{Value: false}
	NULL_VAL  = &NullValue{}
)

type Value interface {
	Type() ValueType
	Inspect() string
}

type IntValue struct {
	Value int32
}

func (v *IntValue) Type() ValueType {
	return INT_VAL_T
}

func (v *IntValue) Inspect() string {
	return strconv.Itoa(int(v.Value))
}

func NewIntValue(val int32) *IntValue {
	return &IntValue{Value: val}
}

type BoolValue struct {
	Value bool
}

func (v *BoolValue) Type() ValueType {
	return BOOL_VAL_T
}

func (v *BoolValue) Inspect() string {
	return strconv.FormatBool(v.Value)
}

func NewBoolValue(val bool) *BoolValue {
	if val {
		return TRUE_VAL
	} else {
		return FALSE_VAL
	}
}

type NullValue struct{}

func (v *NullValue) Type() ValueType {
	return NULL_VAL_T
}

func (v *NullValue) Inspect() string {
	return "null"
}

type ErrorValue struct {
	Message string
}

func (v *ErrorValue) Type() ValueType {
	return ERROR_VAL_T
}

func (v *ErrorValue) Inspect() string {
	return "runtime error: " + v.Message
}

func NewErrorValue(format string, a ...any) *ErrorValue {
	return &ErrorValue{Message: fmt.Sprintf(format, a...)}
}

func IsError(val Value) bool {
	if val != nil {
		return val.Type() == ERROR_VAL_T
	} else {
		return false
	}
}

type FnValue struct {
	Params []*Identifier
	Body   *BlockStatement
	Env    *Environment
}

func (v *FnValue) Type() ValueType {
	return FN_VAL_T
}

func (v *FnValue) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range v.Params {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(v.Body.String())

	return out.String()
}

type ReturnValue struct {
	Value Value
}

func (v *ReturnValue) Type() ValueType {
	return RETURN_VAL_T
}

func (v *ReturnValue) Inspect() string {
	return "return value: " + v.Value.Inspect()
}

func unwrapReturnValue(val Value) Value {
	if returnVal, ok := val.(*ReturnValue); ok {
		return returnVal.Value
	} else {
		return val
	}
}

type StringValue struct {
	Value string
}

func (v *StringValue) Type() ValueType {
	return STRING_VAL_T
}

func (v *StringValue) Inspect() string {
	return v.Value
}

func NewStringValue(val string) *StringValue {
	return &StringValue{Value: val}
}

type BuiltinFnValue struct {
	Fn func(args ...Value) Value
}

func (v *BuiltinFnValue) Type() ValueType {
	return BUILTIN_FN_VAL
}

func (v *BuiltinFnValue) Inspect() string {
	return "builtin fn"
}
