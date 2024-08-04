package main

import (
	"errors"
	"fmt"
)

type Interpreter struct {
	parser     *Parser
	lastResult Value
	lastError  error
}

func NewInterpreter(code string) *Interpreter {
	return &Interpreter{
		parser: NewParser(code),
	}
}

func (interpreter *Interpreter) Eval() (Value, error) {
	ast, err := interpreter.parser.Parse()
	if err != nil {
		return nil, err
	}
	return interpreter.evalAst(ast)
}

func (interpreter *Interpreter) visitNumberExpr(numberExpr *NumberExpr) {
	interpreter.lastResult = NewNumValue(numberExpr.Value)
	interpreter.lastError = nil
}

func (interpreter *Interpreter) visitBooleanExpr(booleanExpr *BooleanExpr) {
	interpreter.lastResult = NewBooleanValue(booleanExpr.Value)
	interpreter.lastError = nil
}

func (interpreter *Interpreter) visitNilExpr() {
	interpreter.lastResult = NewNilValue()
	interpreter.lastError = nil
}

func (interpreter *Interpreter) visitStringExpr(stringExpr *StringExpr) {
	interpreter.lastResult = NewStringValue(stringExpr.Value)
	interpreter.lastError = nil
}

func (interpreter *Interpreter) visitGroupExpr(groupExpr *GroupExpr) {
	interpreter.lastResult, interpreter.lastError = interpreter.evalAst(groupExpr.Inner)
}

func (interpreter *Interpreter) visitUnaryExpr(unaryExpr *UnaryExpr) {
	value, err := interpreter.evalAst(unaryExpr.Value)
	if err != nil {
		interpreter.lastError = nil
		interpreter.lastError = err
		return
	}

	if unaryExpr.Operator.GetLexeme() == "-" {
		if value.getType() == VtNumber {
			num := value.(*NumValue)
			interpreter.lastResult = NewNumValue(-num.Value)
			interpreter.lastError = nil
		} else {
			interpreter.lastResult = nil
			interpreter.lastError = errors.New("unary operator '-' supports only numbers")
		}
	} else {
		switch value.getType() {
		case VtNumber:
			num := value.(*NumValue)
			interpreter.lastResult = NewBooleanValue(num.Value == 0)
			interpreter.lastError = nil
		case VtString:
			str := value.(*StringValue)
			interpreter.lastResult = NewBooleanValue(len(str.Value) == 0)
			interpreter.lastError = nil
		case VtBoolean:
			boolValue := value.(*BooleanValue)
			interpreter.lastResult = NewBooleanValue(!boolValue.Value)
			interpreter.lastError = nil
		case VtNil:
			interpreter.lastResult = NewBooleanValue(true)
			interpreter.lastError = nil
		default:
			interpreter.lastResult = nil
			interpreter.lastError = errors.New("unsupported value type for unary operator '!'")
		}
	}
}

func (interpreter *Interpreter) visitBinaryExpr(expr *BinaryExpr) {
	var left, right Value
	var err error

	interpreter.lastResult = nil
	interpreter.lastError = nil

	left, err = interpreter.evalAst(expr.Left)
	if err != nil {
		interpreter.lastError = err
		return
	}

	right, err = interpreter.evalAst(expr.Right)
	if err != nil {
		interpreter.lastError = err
		return
	}

	leftType := left.getType()
	rightType := right.getType()
	bothNums := leftType == VtNumber && rightType == VtNumber

	switch op := expr.Operator.GetLexeme(); op {
	case "*", "/", "-": //, ">", ">=", "<", "<=":
		if bothNums {
			leftNum := left.(*NumValue).Value
			rightNum := right.(*NumValue).Value
			switch op {
			case "*":
				interpreter.lastResult = NewNumValue(leftNum * rightNum)
			case "/":
				interpreter.lastResult = NewNumValue(leftNum / rightNum)
			case "-":
				interpreter.lastResult = NewNumValue(leftNum - rightNum)
			}
		} else {
			interpreter.lastError = errors.New("only numbers are supported")
		}
	case "+":
		if bothNums {
			interpreter.lastResult = NewNumValue(left.(*NumValue).Value + right.(*NumValue).Value)
		} else if leftType == VtString && rightType == VtString {
			interpreter.lastResult = NewStringValue(left.(*StringValue).Value + right.(*StringValue).Value)
		}
	default:
		interpreter.lastError = errors.New(fmt.Sprintf("unsupported operator %s", op))
	}

}

func (interpreter *Interpreter) evalAst(ast AST) (Value, error) {
	ast.accept(interpreter)
	return interpreter.lastResult, interpreter.lastError
}
