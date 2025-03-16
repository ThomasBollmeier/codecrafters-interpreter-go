package main

import (
	"fmt"
)

type varInfo struct {
	parent *varInfo
	vars   map[string]int
}

func newVarInfo(parent *varInfo) *varInfo {
	return &varInfo{
		parent: parent,
		vars:   make(map[string]int),
	}
}

func (v *varInfo) getLevel(name string) (int, error) {
	_, ok := v.vars[name]
	if ok {
		return 0, nil
	}
	if v.parent != nil {
		level, err := v.parent.getLevel(name)
		if err != nil {
			return -1, err
		}
		return level + 1, nil
	}
	return -1, fmt.Errorf("variable %s not found", name)
}

func (v *varInfo) addName(name string) {
	v.vars[name] = 1
}

type VariableResolver struct {
	varInfo *varInfo
}

func NewVariableResolver() *VariableResolver {
	return &VariableResolver{
		varInfo: newVarInfo(nil),
	}
}

func (v *VariableResolver) visitProgram(program *Program) {
	for _, stmt := range program.statements {
		stmt.accept(v)
	}
}

func (v *VariableResolver) visitBlock(block *Block) {
	v.varInfo = newVarInfo(v.varInfo)
	for _, stmt := range block.statements {
		stmt.accept(v)
	}
	v.varInfo = v.varInfo.parent
}

func (v *VariableResolver) visitVarDecl(varDecl *VarDecl) {
	varDecl.expression.accept(v)
	v.varInfo.addName(varDecl.name)
}

func (v *VariableResolver) visitPrint(printStmt *PrintStatement) {
	printStmt.expression.accept(v)
}

func (v *VariableResolver) visitReturnStmt(returnStmt *ReturnStatement) {
	if returnStmt.expression != nil {
		returnStmt.expression.accept(v)
	}
}

func (v *VariableResolver) visitExprStmt(exprStmt *ExpressionStatement) {
	exprStmt.expression.accept(v)
}

func (v *VariableResolver) visitIfStmt(ifStmt *IfStatement) {
	ifStmt.condition.accept(v)
	ifStmt.consequent.accept(v)
	if ifStmt.alternate != nil {
		ifStmt.alternate.accept(v)
	}
}

func (v *VariableResolver) visitWhileStmt(whileStmt *WhileStatement) {
	whileStmt.condition.accept(v)
	whileStmt.statement.accept(v)
}

func (v *VariableResolver) visitForStmt(f *ForStatement) {
	if f.initializer != nil {
		f.initializer.accept(v)
	}
	if f.condition != nil {
		f.condition.accept(v)
	}
	if f.increment != nil {
		f.increment.accept(v)
	}
	f.statement.accept(v)
}

func (v *VariableResolver) visitFunctionDef(f *FunctionDef) {
	v.varInfo.addName(f.name)
	v.varInfo = newVarInfo(v.varInfo)
	for _, param := range f.parameters {
		v.varInfo.addName(param)
	}
	f.body.accept(v)
	v.varInfo = v.varInfo.parent
}

func (v *VariableResolver) visitNumberExpr(*NumberExpr) {}

func (v *VariableResolver) visitBooleanExpr(*BooleanExpr) {}

func (v *VariableResolver) visitNilExpr() {}

func (v *VariableResolver) visitStringExpr(*StringExpr) {}

func (v *VariableResolver) visitIdentifierExpr(identifierExpr *IdentifierExpr) {
	identifierExpr.defLevel, _ = v.varInfo.getLevel(identifierExpr.name)
}

func (v *VariableResolver) visitGroupExpr(groupExpr *GroupExpr) {
	groupExpr.Inner.accept(v)
}

func (v *VariableResolver) visitUnaryExpr(unaryExpr *UnaryExpr) {
	unaryExpr.Value.accept(v)
}

func (v *VariableResolver) visitBinaryExpr(expr *BinaryExpr) {
	expr.Left.accept(v)
	expr.Right.accept(v)
}

func (v *VariableResolver) visitAssignment(assignment *Assignment) {
	assignment.right.accept(v)
	assignment.defLevel, _ = v.varInfo.getLevel(assignment.left)
}

func (v *VariableResolver) visitCall(call *Call) {
	call.callee.accept(v)
	for _, arg := range call.args {
		arg.accept(v)
	}
}
