package main

import (
	"errors"
	"fmt"
)

type Interpreter struct {
	parser     *Parser
	lastResult Value
	lastError  error
	env        Environment
}

func NewInterpreter(code string) *Interpreter {
	return &Interpreter{
		parser: NewParser(code),
		env:    *NewEnvironment(nil),
	}
}

func NewInterpreterWithEnv(code string, env Environment) *Interpreter {
	return &Interpreter{
		parser: NewParser(code),
		env:    env,
	}
}

func (interpreter *Interpreter) Run() (error, bool) {
	ast, err := interpreter.parser.ParseProgram()
	if err != nil {
		return err, false
	}
	ast.accept(interpreter)
	return interpreter.lastError, interpreter.lastError != nil
}

func (interpreter *Interpreter) Eval() (Value, error, bool) {
	ast, err := interpreter.parser.ParseExpression()
	if err != nil {
		return nil, err, false
	}
	var value Value
	value, err = interpreter.evalAst(ast)
	return value, err, err != nil
}

func (interpreter *Interpreter) visitProgram(program *Program) {
	for _, statement := range program.statements {
		statement.accept(interpreter)
		if interpreter.lastError != nil {
			break
		}
	}
}

func (interpreter *Interpreter) visitVarDecl(varDecl *VarDecl) {
	value, err := interpreter.evalAst(varDecl.expression)
	if err != nil {
		return
	}
	interpreter.env.Set(varDecl.name, value)
}

func (interpreter *Interpreter) visitPrint(printStmt *PrintStatement) {
	value, err := interpreter.evalAst(printStmt.expression)
	if err == nil {
		fmt.Println(value)
	}
}

func (interpreter *Interpreter) visitExprStmt(exprStmt *ExpressionStatement) {
	_, _ = interpreter.evalAst(exprStmt.expression)
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

func (interpreter *Interpreter) visitIdentifierExpr(identifierExpr *IdentifierExpr) {
	interpreter.lastResult, interpreter.lastError = interpreter.env.Get(identifierExpr.name)
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
	case "*", "/", "-", ">", ">=", "<", "<=":
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
			case ">":
				interpreter.lastResult = NewBooleanValue(leftNum > rightNum)
			case ">=":
				interpreter.lastResult = NewBooleanValue(leftNum >= rightNum)
			case "<":
				interpreter.lastResult = NewBooleanValue(leftNum < rightNum)
			case "<=":
				interpreter.lastResult = NewBooleanValue(leftNum <= rightNum)
			}
		} else {
			interpreter.lastError = errors.New("only numbers are supported")
		}
	case "+":
		if bothNums {
			interpreter.lastResult = NewNumValue(left.(*NumValue).Value + right.(*NumValue).Value)
		} else if leftType == VtString && rightType == VtString {
			interpreter.lastResult = NewStringValue(left.(*StringValue).Value + right.(*StringValue).Value)
		} else {
			interpreter.lastError = errors.New("only two numbers or two strings are supported")
		}
	case "==":
		interpreter.lastResult = NewBooleanValue(left.isEqualTo(right))
	case "!=":
		interpreter.lastResult = NewBooleanValue(!left.isEqualTo(right))
	default:
		interpreter.lastError = errors.New(fmt.Sprintf("unsupported operator %s", op))
	}

}

func (interpreter *Interpreter) visitAssignment(assignment *Assignment) {
	value, err := interpreter.evalAst(assignment.right)
	if err != nil {
		return
	}
	interpreter.env.Set(assignment.left, value)
}

func (interpreter *Interpreter) evalAst(ast AST) (Value, error) {
	ast.accept(interpreter)
	return interpreter.lastResult, interpreter.lastError
}
