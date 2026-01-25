package main

import (
	"bytes"
	"fmt"
)

var builtins = map[string]*BuiltinFnValue{
	"echo": &BuiltinFnValue{
		Fn: func(args ...Value) Value {
			if len(args) != 1 {
				return NewErrorValue("wrong number of arguments: wanted 1, got %d", len(args))
			}

			return args[0]
		},
	},
	"inspect": &BuiltinFnValue{
		Fn: func(args ...Value) Value {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL_VAL
		},
	},
	"println": &BuiltinFnValue{
		Fn: func(args ...Value) Value {
			if len(args) < 1 || args[0].Type() != STRING_VAL_T {
				return NewErrorValue("println expects at least one argument of type string")
			}

			input := args[0].(*StringValue).Value
			argInd := 1
			var res bytes.Buffer
			for i := 0; i < len(input); i++ {
				if input[i] == '{' && input[i+1] == '}' {
					if argInd >= len(args) {
						return NewErrorValue("println expects more arguments")
					}
					i++
					res.WriteString(args[argInd].Inspect())
					argInd++
				} else {
					res.WriteByte(input[i])
				}
			}

			fmt.Println(res.String())
			return NULL_VAL
		},
	},
}

func Eval(node Node, env *Environment) Value {
	switch node := node.(type) {
	case *Program:
		return unwrapReturnValue(evalStatements(node.Statements, env))
	case *LetStatement:
		return evalLetStatement(node, env)
	case *BlockStatement:
		return evalStatements(node.Statements, env.NewChild())
	case *IfStatement:
		return evalIfStatement(node, env)
	case *ReturnStatement:
		return &ReturnValue{Value: Eval(node.Value, env)}
	case *FnStatement:
		stmt := &LetStatement{Name: node.Name, Value: node.Value}
		return evalLetStatement(stmt, env)
	case *ExprStatement:
		return Eval(node.Value, env)
	case *IntLiteral:
		return &IntValue{Value: node.Value}
	case *BoolLiteral:
		return NewBoolValue(node.Value)
	case *NullLiteral:
		return &NullValue{}
	case *Identifier:
		return evalIdentifier(node, env)
	case *FnLiteral:
		return &FnValue{Params: node.Params, Body: node.Body, Env: env}
	case *StringLiteral:
		return &StringValue{Value: node.Value}
	case *PrefixExpression:
		right := Eval(node.Right, env)
		if IsError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *InfixExpression:
		right := Eval(node.Right, env)
		if IsError(right) {
			return right
		}

		if node.Operator == "=" {
			return evalAssignmentExpression(node.Left, right, env)
		}

		left := Eval(node.Left, env)
		if IsError(left) {
			return left
		}

		return evalInfixExpression(node.Operator, left, right)
	case *CallExpression:
		fn := Eval(node.Function, env)
		if IsError(fn) {
			return fn
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && IsError(args[0]) {
			return args[0]
		}

		return evalCall(fn, args)
	default:
		return nil
	}
}

func evalLetStatement(stmt *LetStatement, env *Environment) Value {
	val := Eval(stmt.Value, env)
	if IsError(val) {
		return val
	}

	if env.Let(stmt.Name.Value, val) {
		return NULL_VAL
	} else {
		return NewErrorValue("identifier already exists: %s", stmt.Name.Value)
	}
}

func evalStatements(stmts []Statement, env *Environment) Value {
	var res Value

	for _, stmt := range stmts {
		res = Eval(stmt, env)

		switch res := res.(type) {
		case *ReturnValue:
			return res
		case *ErrorValue:
			return res
		}
	}

	return res
}

func evalIfStatement(stmt *IfStatement, env *Environment) Value {
	condition := Eval(stmt.Condition, env)
	if condition.Type() != BOOL_VAL_T {
		return NewErrorValue("if statement condition must evaluate to bool not %s", condition.Type())
	}

	if condition.(*BoolValue).Value {
		return Eval(stmt.Then, env)
	} else if stmt.Else != nil {
		return Eval(stmt.Else, env)
	} else {
		return NULL_VAL
	}
}

func evalExpressions(exprs []Expression, env *Environment) []Value {
	var res []Value

	for _, expr := range exprs {
		val := Eval(expr, env)
		if IsError(val) {
			return []Value{val}
		}

		res = append(res, val)
	}

	return res
}

func evalCall(fn Value, args []Value) Value {
	switch fn := fn.(type) {
	case *FnValue:
		if len(args) != len(fn.Params) {
			return NewErrorValue("wrong number of arguments: wanted %d, got %d", len(fn.Params), len(args))
		}

		env := extendFnEnv(fn, args)
		val := evalStatements(fn.Body.Statements, env)
		return unwrapReturnValue(val)
	case *BuiltinFnValue:
		return fn.Fn(args...)
	default:
		return NewErrorValue("%s is not a function", fn.Type())
	}
}

func extendFnEnv(fn *FnValue, args []Value) *Environment {
	env := fn.Env.NewChild()

	for paramIndex, param := range fn.Params {
		env.Let(param.Value, args[paramIndex])
	}

	return env
}

func evalIdentifier(ident *Identifier, env *Environment) Value {
	if val, ok := env.Get(ident.Value); ok {
		return val
	}

	if val, ok := builtins[ident.Value]; ok {
		return val
	}

	return NewErrorValue("identifier not found: %s", ident.Value)
}

func evalPrefixExpression(operator string, right Value) Value {
	switch operator {
	case "!":
		return evalBangExpression(right)
	case "-":
		return evalMinusPrefixExpression(right)
	default:
		return nil
	}
}

func evalBangExpression(right Value) Value {
	switch right {
	case TRUE_VAL:
		return FALSE_VAL
	case FALSE_VAL:
		return TRUE_VAL
	default:
		return NewErrorValue("type %s is incompatible with ! operator", right.Type())
	}
}

func evalMinusPrefixExpression(right Value) Value {
	if right.Type() != INT_VAL_T {
		return NewErrorValue("type %s is incompatible with - operator", right.Type())
	}

	val := right.(*IntValue).Value
	return &IntValue{Value: -val}
}

func evalInfixExpression(operator string, left, right Value) Value {
	if left.Type() == INT_VAL_T && right.Type() == INT_VAL_T {
		return evalIntInfixExpression(operator, left.(*IntValue), right.(*IntValue))
	} else if left.Type() == BOOL_VAL_T && right.Type() == BOOL_VAL_T {
		return evalBoolInfixExpression(operator, left.(*BoolValue), right.(*BoolValue))
	} else if left.Type() == STRING_VAL_T && right.Type() == STRING_VAL_T {
		return evalStringInfixExpression(operator, left.(*StringValue), right.(*StringValue))
	} else {
		return NewErrorValue("types mismatch %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntInfixExpression(operator string, left, right *IntValue) Value {
	leftVal := left.Value
	rightVal := right.Value

	switch operator {
	case "+":
		return NewIntValue(leftVal + rightVal)
	case "-":
		return NewIntValue(leftVal - rightVal)
	case "*":
		return NewIntValue(leftVal * rightVal)
	case "/":
		return NewIntValue(leftVal / rightVal)
	case "==":
		return NewBoolValue(leftVal == rightVal)
	case "!=":
		return NewBoolValue(leftVal != rightVal)
	case "<":
		return NewBoolValue(leftVal < rightVal)
	case "<=":
		return NewBoolValue(leftVal <= rightVal)
	case ">":
		return NewBoolValue(leftVal > rightVal)
	case ">=":
		return NewBoolValue(leftVal >= rightVal)
	default:
		return NewErrorValue("type int is incompatible with %s operator", operator)
	}
}

func evalBoolInfixExpression(operator string, left, right *BoolValue) Value {
	switch operator {
	case "==":
		return NewBoolValue(left == right)
	case "!=":
		return NewBoolValue(left != right)
	default:
		return NewErrorValue("type bool is incompatible with %s operator", operator)
	}
}

func evalStringInfixExpression(operator string, left, right *StringValue) Value {
	switch operator {
	case "+":
		return NewStringValue(left.Value + right.Value)
	default:
		return NewErrorValue("type string is incompatible with %s operator", operator)
	}
}

func evalAssignmentExpression(left Expression, right Value, env *Environment) Value {
	switch ident := left.(type) {
	case *Identifier:
		ok := env.Assign(ident.Value, right)
		if ok {
			return right
		} else {
			return NewErrorValue("identifier not found: %s", ident.Value)
		}
	default:
		return NewErrorValue("%s is not an identifier", ident.String())
	}
}
