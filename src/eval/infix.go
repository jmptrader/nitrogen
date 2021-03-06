package eval

import (
	"math"

	"github.com/nitrogen-lang/nitrogen/src/object"
)

func (i *Interpreter) evalInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return object.NewException("type mismatch: %s %s %s", left.Type(), op, right.Type())
	case object.ObjectsAre(object.IntergerObj, left, right):
		return i.evalIntegerInfixExpression(op, left, right)
	case object.ObjectsAre(object.FloatObj, left, right):
		return i.evalFloatInfixExpression(op, left, right)
	case object.ObjectsAre(object.StringObj, left, right):
		return i.evalStringInfixExpression(op, left, right)
	case object.ObjectsAre(object.ArrayObj, left, right):
		return i.evalArrayInfixExpression(op, left, right)
	case object.ObjectsAre(object.BooleanObj, left, right):
		return i.evalBoolInfixExpression(op, left, right)
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func (i *Interpreter) evalIntegerInfixExpression(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return object.NativeBoolToBooleanObj(leftVal < rightVal)
	case ">":
		return object.NativeBoolToBooleanObj(leftVal > rightVal)
	case "==":
		return object.NativeBoolToBooleanObj(leftVal == rightVal)
	case "!=":
		return object.NativeBoolToBooleanObj(leftVal != rightVal)
	case "<=":
		return object.NativeBoolToBooleanObj(leftVal <= rightVal)
	case ">=":
		return object.NativeBoolToBooleanObj(leftVal >= rightVal)
	case "<<":
		if rightVal < 0 {
			return object.NewException("Shift value must be non-negative")
		}
		return &object.Integer{Value: leftVal << uint64(rightVal)}
	case ">>":
		if rightVal < 0 {
			return object.NewException("Shift value must be non-negative")
		}
		return &object.Integer{Value: leftVal >> uint64(rightVal)}
	case "&":
		return &object.Integer{Value: leftVal & rightVal}
	case "&^":
		return &object.Integer{Value: leftVal &^ rightVal}
	case "|":
		return &object.Integer{Value: leftVal | rightVal}
	case "^":
		return &object.Integer{Value: leftVal ^ rightVal}
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func (i *Interpreter) evalFloatInfixExpression(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch op {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "%":
		return &object.Float{Value: math.Mod(leftVal, rightVal)}
	case "<":
		return object.NativeBoolToBooleanObj(leftVal < rightVal)
	case ">":
		return object.NativeBoolToBooleanObj(leftVal > rightVal)
	case "==":
		return object.NativeBoolToBooleanObj(leftVal == rightVal)
	case "!=":
		return object.NativeBoolToBooleanObj(leftVal != rightVal)
	case "<=":
		return object.NativeBoolToBooleanObj(leftVal <= rightVal)
	case ">=":
		return object.NativeBoolToBooleanObj(leftVal >= rightVal)
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func (i *Interpreter) evalStringInfixExpression(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch op {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return object.NativeBoolToBooleanObj(leftVal == rightVal)
	case "!=":
		return object.NativeBoolToBooleanObj(leftVal != rightVal)
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func (i *Interpreter) evalArrayInfixExpression(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.Array)
	rightVal := right.(*object.Array)

	switch op {
	case "+":
		leftLen := len(leftVal.Elements)
		rightLen := len(rightVal.Elements)
		newElements := make([]object.Object, leftLen+rightLen, leftLen+rightLen)
		copy(newElements, leftVal.Elements)
		copy(newElements[leftLen:], rightVal.Elements)
		return &object.Array{Elements: newElements}
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func (i *Interpreter) evalBoolInfixExpression(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch op {
	case "==":
		return object.NativeBoolToBooleanObj(leftVal == rightVal)
	case "!=":
		return object.NativeBoolToBooleanObj(leftVal != rightVal)
	case "or":
		return object.NativeBoolToBooleanObj(leftVal || rightVal)
	case "and":
		return object.NativeBoolToBooleanObj(leftVal && rightVal)
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}
