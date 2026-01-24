package main

import (
	"fmt"
	"strconv"
)

type ValueType string

const (
	INT_VAL_T   = "int"
	BOOL_VAL_T  = "bool"
	NULL_VAL_T  = "null"
	ERROR_VAL_T = "error"
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
