package main

func Eval(node Node, env *Environment) Value {
	switch node := node.(type) {
	case *Program:
		return evalProgram(node, env)
	case *LetStatement:
		return evalLetStatement(node, env)
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
			switch ident := node.Left.(type) {
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

		left := Eval(node.Left, env)
		if IsError(left) {
			return left
		}

		return evalInfixExpression(node.Operator, left, right)
	default:
		return nil
	}
}

func evalProgram(program *Program, env *Environment) Value {
	var res Value

	for _, stmt := range program.Statements {
		res = Eval(stmt, env)

		switch res := res.(type) {
		case *ErrorValue:
			return res
		}
	}

	return res
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

func evalIdentifier(ident *Identifier, env *Environment) Value {
	val, ok := env.Get(ident.Value)
	if !ok {
		return NewErrorValue("identifier not found: %s", ident.Value)
	}

	return val
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
