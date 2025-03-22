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
	interpreter.env.StartDeclaration(varDecl.name) // start declaration
	value, err := interpreter.evalAst(varDecl.expression)
	if err != nil {
		return
	}
	interpreter.env.Set(varDecl.name, value) // finalize declaration
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

	interpreter.env = NewEnvironment(interpreter.env)
	defer func() {
		interpreter.env = interpreter.env.parent
	}()

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

func (interpreter *Interpreter) visitClassDef(c *ClassDef) {
	var methods []LambdaValue
	var method Value
	var err error
	interpreter.env = NewEnvironment(interpreter.env)

	for _, function := range c.functions {
		method, err = interpreter.evalAst(&function)
		if err != nil {
			break
		}
		methods = append(methods, *method.(*LambdaValue))
	}

	interpreter.env = interpreter.env.parent

	if err == nil {
		class := NewClassValue(c.name, methods)
		interpreter.env.Set(c.name, class)
		interpreter.lastResult = class
		interpreter.lastError = nil
	}
}

func (interpreter *Interpreter) visitFunctionDef(funDef *FunctionDef) {
	var name string
	isConstructor := false
	if funDef.class == nil {
		name = funDef.name
	} else {
		name = funDef.class.name + "::" + funDef.name
		isConstructor = funDef.name == "init"
	}
	lambda := NewLambdaValue(name, funDef.parameters, funDef.body, *interpreter.env)
	lambda.isConstructor = isConstructor
	interpreter.env.Set(name, lambda)
	interpreter.lastResult = lambda
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
	case ".":
		interpreter.lastResult, interpreter.lastError = interpreter.evalPathExpr(expr)
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

	identifier, isIdent := assignment.left.(*IdentifierExpr)
	if isIdent {
		var defEnv *Environment
		if assignment.defLevel != -1 {
			defEnv, err = interpreter.env.GetEnvAtLevel(assignment.defLevel)
		} else {
			defEnv, err = interpreter.env.GetDefiningEnv(identifier.name)
		}
		if err != nil {
			return
		}
		defEnv.Set(identifier.name, value)
		return
	}

	pathExpr := assignment.left.(*BinaryExpr)
	instance, property, err := interpreter.evalPathExprLhs(pathExpr)
	if err != nil {
		interpreter.lastResult = nil
		interpreter.lastError = err
		return
	}
	err = instance.setProperty(property, value)
	if err != nil {
		interpreter.lastResult = nil
		interpreter.lastError = err
		return
	}

	interpreter.lastResult = NewNilValue()
	interpreter.lastError = nil
}

func (interpreter *Interpreter) visitCall(call *Call) {
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

	arguments, err := interpreter.evalArguments(call.args)
	if err != nil {
		return
	}

	interpreter.lastResult, interpreter.lastError = callableValue.call(arguments)
}

func (interpreter *Interpreter) evalArguments(args []Expr) ([]Value, error) {
	var arguments []Value

	for _, arg := range args {
		argument, err := interpreter.evalAst(arg)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)
	}
	return arguments, nil
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

func (interpreter *Interpreter) evalPathExprLhs(expr *BinaryExpr) (*InstanceValue, string, error) {
	value, err := interpreter.evalAst(expr.Left)
	if err != nil {
		return nil, "", err
	}
	instance, isInstance := value.(*InstanceValue)
	if !isInstance {
		return nil, "", fmt.Errorf("expected instance but got %T", value)
	}
	return interpreter.evalPathLhs(instance, expr.Right)
}

func (interpreter *Interpreter) evalPathLhs(instance *InstanceValue, expr Expr) (*InstanceValue, string, error) {
	ident, isIdent := expr.(*IdentifierExpr)
	if isIdent {
		return instance, ident.name, nil
	}

	binExpr, isBinExpr := expr.(*BinaryExpr)
	if isBinExpr {
		next, err := interpreter.evalPath(instance, binExpr.Left)
		if err != nil {
			return nil, "", err
		}
		nextInstance, isInstance := next.(*InstanceValue)
		if !isInstance {
			return nil, "", errors.New("expected expression to evaluate to an instance")
		}
		return interpreter.evalPathLhs(nextInstance, binExpr.Right)
	}

	return nil, "", errors.New("invalid expression as lhs")
}

func (interpreter *Interpreter) evalPathExpr(expr *BinaryExpr) (Value, error) {
	value, err := interpreter.evalAst(expr.Left)
	if err != nil {
		return nil, err
	}
	instance, isInstance := value.(*InstanceValue)
	if !isInstance {
		return nil, fmt.Errorf("expected instance but got %T", value)
	}
	return interpreter.evalPath(instance, expr.Right)
}

func (interpreter *Interpreter) evalPath(instance *InstanceValue, expr Expr) (Value, error) {
	ident, isIdent := expr.(*IdentifierExpr)
	if isIdent {
		return instance.getMember(ident.name)
	}

	call, isCall := expr.(*Call)
	if isCall {
		method, errMethod := interpreter.evalMethod(instance, call.callee)
		if errMethod != nil {
			return nil, errMethod
		}

		arguments, errArgs := interpreter.evalArguments(call.args)
		if errArgs != nil {
			return nil, errArgs
		}
		return method.call(arguments)
	}

	binExpr, isBinExpr := expr.(*BinaryExpr)
	if isBinExpr {
		next, err := interpreter.evalPath(instance, binExpr.Left)
		if err != nil {
			return nil, err
		}
		nextInstance, isInstance := next.(*InstanceValue)
		if !isInstance {
			return nil, errors.New("expected expression to evaluate to an instance")
		}
		return interpreter.evalPath(nextInstance, binExpr.Right)
	}

	return nil, errors.New("invalid path segment")
}

func (interpreter *Interpreter) evalMethod(instance *InstanceValue, callee Expr) (callable, error) {
	ident, isIdent := callee.(*IdentifierExpr)
	if isIdent {
		member, errMember := instance.getMember(ident.name)
		if errMember != nil {
			return nil, errMember
		}
		method, isMethod := member.(callable)
		if !isMethod {
			return nil, errors.New(fmt.Sprintf("expected callable method but got %T", member))
		}
		return method, nil
	}

	call, isCall := callee.(*Call)
	if isCall {
		calleeValue, errCallee := interpreter.evalMethod(instance, call.callee)
		if errCallee != nil {
			return nil, errCallee
		}
		arguments, errArgs := interpreter.evalArguments(call.args)
		if errArgs != nil {
			return nil, errArgs
		}
		value, errCall := calleeValue.call(arguments)
		if errCall != nil {
			return nil, errCall
		}
		method, isMethod := value.(callable)
		if !isMethod {
			return nil, errors.New(fmt.Sprintf("expected callable method but got %T", value))
		}
		return method, nil
	}

	return nil, fmt.Errorf("invalid callee type: %T", callee)
}

func (interpreter *Interpreter) evalAst(ast AST) (Value, error) {
	ast.accept(interpreter)
	return interpreter.lastResult, interpreter.lastError
}
