package main

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
	ast.accept(interpreter)
	return interpreter.lastResult, interpreter.lastError
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
	panic("not implemented")
}

func (interpreter *Interpreter) visitGroupExpr(groupExpr *GroupExpr) {
	panic("not implemented")
}

func (interpreter *Interpreter) visitUnaryExpr(unaryExpr *UnaryExpr) {
	panic("not implemented")
}

func (interpreter *Interpreter) visitBinaryExpr(expr *BinaryExpr) {
	panic("not implemented")
}
