package main

import (
	"errors"
	"fmt"
)

type Interpreter struct {
	lastResult       Value
	lastError        error
	lambdaEvalActive bool
	returnOccurred   bool
	env              *Environment
}

func NewInterpreter(env *Environment) *Interpreter {
	if env == nil {
		return &Interpreter{
			env: NewEnvironment(nil),
		}
	} else {
		return &Interpreter{
			env: env,
		}
	}
}

func (interpreter *Interpreter) Run(code string) (error, bool) {
	parser := NewParser(code)
	ast, err := parser.ParseProgram()
	if err != nil {
		return err, false
	}
	ast.accept(interpreter)
	return interpreter.lastError, interpreter.lastError != nil
}

func (interpreter *Interpreter) Eval(code string) (Value, error, bool) {
	parser := NewParser(code)
	ast, err := parser.ParseExpression()
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

func (interpreter *Interpreter) visitBlock(block *Block) {
	blockEnv := NewEnvironment(interpreter.env)
	interpreter.env = blockEnv

	for _, statement := range block.statements {
		statement.accept(interpreter)
		if interpreter.lastError != nil || interpreter.returnOccurred {
			break
		}
	}

	interpreter.env = blockEnv.parent
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

func (interpreter *Interpreter) visitReturnStmt(returnStmt *ReturnStatement) {
	if !interpreter.lambdaEvalActive {
		interpreter.lastResult = nil
		interpreter.lastError = errors.New("return is not allowed outside of a function body")
		return
	}
	if returnStmt.expression == nil {
		interpreter.lastResult = NewNilValue()
		interpreter.lastError = nil
	} else {
		interpreter.lastResult, interpreter.lastError = interpreter.evalAst(returnStmt.expression)
	}

	interpreter.returnOccurred = true
}

func (interpreter *Interpreter) visitExprStmt(exprStmt *ExpressionStatement) {
	_, _ = interpreter.evalAst(exprStmt.expression)
}

func (interpreter *Interpreter) visitIfStmt(ifStmt *IfStatement) {
	value, err := interpreter.evalAst(ifStmt.condition)
	if err != nil {
		return
	}
	if value.isTruthy() {
		value, err = interpreter.evalAst(ifStmt.consequent)
		if err != nil {
			return
		}
		interpreter.lastResult = value
		interpreter.lastError = nil
	} else if ifStmt.alternate != nil {
		value, err = interpreter.evalAst(ifStmt.alternate)
		if err != nil {
			return
		}
		interpreter.lastResult = value
		interpreter.lastError = nil
	} else {
		interpreter.lastResult = NewNilValue()
		interpreter.lastError = nil
	}
}

func (interpreter *Interpreter) visitWhileStmt(whileStmt *WhileStatement) {
	for {
		condition, err := interpreter.evalAst(whileStmt.condition)
		if err != nil {
			return
		}
		if !condition.isTruthy() {
			break
		}

		_, err = interpreter.evalAst(whileStmt.statement)
		if err != nil || interpreter.returnOccurred {
			return
		}
	}
	interpreter.lastResult = NewNilValue()
	interpreter.lastError = nil
}

func (interpreter *Interpreter) visitForStmt(forStmt *ForStatement) {
	var err error

	if forStmt.initializer != nil {
		_, err = interpreter.evalAst(forStmt.initializer)
		if err != nil {
			return
		}
	}

	var condVal Value

	for {
		if forStmt.condition != nil {
			condVal, err = interpreter.evalAst(forStmt.condition)
			if err != nil {
				return
			}
			if !condVal.isTruthy() {
				break
			}
		}

		_, err = interpreter.evalAst(forStmt.statement)
		if err != nil || interpreter.returnOccurred {
			return
		}

		if forStmt.increment != nil {
			_, err = interpreter.evalAst(forStmt.increment)
			if err != nil {
				return
			}
		}
	}

	interpreter.lastResult = NewNilValue()
	interpreter.lastError = nil
}

func (interpreter *Interpreter) visitFunctionDef(funDef *FunctionDef) {
	lambda := NewLambdaValue(funDef.name, funDef.parameters, funDef.body, *interpreter.env)
	interpreter.env.Set(funDef.name, lambda)
	interpreter.lastResult = NewNilValue()
	interpreter.lastError = nil
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
	var env *Environment
	if identifierExpr.defLevel != -1 {
		env, interpreter.lastError = interpreter.env.GetEnvAtLevel(identifierExpr.defLevel)
		if interpreter.lastError != nil {
			interpreter.lastResult = nil
			return
		}
	} else {
		env = interpreter.env
	}
	interpreter.lastResult, interpreter.lastError = env.Get(identifierExpr.name)
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

	op := expr.Operator.GetLexeme()

	switch op {
	case "and":
		interpreter.lastResult, interpreter.lastError = interpreter.evalConjunction(expr)
		return
	case "or":
		interpreter.lastResult, interpreter.lastError = interpreter.evalDisjunction(expr)
		return
	}

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

	switch op {
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
	var defEnv *Environment
	if assignment.defLevel != -1 {
		defEnv, err = interpreter.env.GetEnvAtLevel(assignment.defLevel)
	} else {
		defEnv, err = interpreter.env.GetDefiningEnv(assignment.left)
	}
	if err != nil {
		return
	}
	defEnv.Set(assignment.left, value)
}

func (interpreter *Interpreter) visitCall(call *Call) {
	var arguments []Value

	for _, arg := range call.args {
		argument, err := interpreter.evalAst(arg)
		if err != nil {
			return
		}
		arguments = append(arguments, argument)
	}

	value, err := interpreter.evalAst(call.callee)
	if err != nil {
		return
	}
	callableValue, ok := value.(callable)
	if !ok {
		interpreter.lastResult = nil
		interpreter.lastError = errors.New("invalid callable")
		return
	}

	interpreter.lastResult, interpreter.lastError = callableValue.call(arguments)
}

func (interpreter *Interpreter) evalDisjunction(expr *BinaryExpr) (Value, error) {
	left, err := interpreter.evalAst(expr.Left)
	if err != nil {
		return nil, err
	}
	if left.isTruthy() {
		return left, nil
	}
	right, err := interpreter.evalAst(expr.Right)
	if err != nil {
		return nil, err
	}
	return right, nil
}

func (interpreter *Interpreter) evalConjunction(expr *BinaryExpr) (Value, error) {
	left, err := interpreter.evalAst(expr.Left)
	if err != nil {
		return nil, err
	}
	if !left.isTruthy() {
		return left, nil
	}
	right, err := interpreter.evalAst(expr.Right)
	if err != nil {
		return nil, err
	}
	return right, nil
}

func (interpreter *Interpreter) evalAst(ast AST) (Value, error) {
	ast.accept(interpreter)
	return interpreter.lastResult, interpreter.lastError
}
