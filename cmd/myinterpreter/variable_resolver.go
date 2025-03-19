package main

import (
	"fmt"
)

type varInfo struct {
	parent          *varInfo
	vars            map[string]int
	isParameterInfo bool
}

func newVarInfo(parent *varInfo) *varInfo {
	return &varInfo{
		parent:          parent,
		vars:            make(map[string]int),
		isParameterInfo: false,
	}
}

func (v *varInfo) getLevel(name string) (int, error) {
	level, ok := v.vars[name]
	if ok {
		if level == 1 {
			return 0, nil
		} else {
			return -1, fmt.Errorf("variable %s is not finally declared", name)
		}
	}
	var err error
	if v.parent != nil {
		level, err = v.parent.getLevel(name)
		if err != nil || level == -1 {
			return -1, err
		}
		return level + 1, nil
	}
	return -1, nil
}

func (v *varInfo) addName(name string) error {
	_, exists := v.vars[name]
	if exists {
		return fmt.Errorf("variable %s already defined", name)
	}
	v.vars[name] = 1
	return nil
}

func (v *varInfo) startVarDecl(name string) error {
	_, exists := v.vars[name]
	if !exists {
		v.vars[name] = -1
		return nil
	}
	return fmt.Errorf("variable %s already declared", name)
}

func (v *varInfo) endVarDecl(name string) {
	v.vars[name] = 1
}

type VariableResolver struct {
	varInfo *varInfo
	err     error
}

func NewVariableResolver() *VariableResolver {
	return &VariableResolver{
		varInfo: newVarInfo(nil),
	}
}

func (v *VariableResolver) visitProgram(program *Program) {
	for _, stmt := range program.statements {
		stmt.accept(v)
		if v.err != nil {
			return
		}
	}
}

func (v *VariableResolver) visitBlock(block *Block) {
	v.varInfo = newVarInfo(v.varInfo)
	for _, stmt := range block.statements {
		stmt.accept(v)
		if v.err != nil {
			break
		}
	}
	v.varInfo = v.varInfo.parent
}

func (v *VariableResolver) visitVarDecl(varDecl *VarDecl) {
	err := v.varInfo.startVarDecl(varDecl.name)
	if err != nil {
		// Is it a self definition like var a = a; ?
		idExpr, ok := varDecl.expression.(*IdentifierExpr)
		isSelfDefinition := ok && idExpr.name == varDecl.name

		// Is it a redeclaration at global scope?
		isGlobalScope := v.varInfo.parent == nil

		if !isSelfDefinition && !isGlobalScope {
			v.err = err
			return
		}
	}
	// Check if a parameter of the same name exists:
	parameterInfo := v.varInfo.parent
	if parameterInfo != nil && parameterInfo.isParameterInfo {
		level, errLevel := parameterInfo.getLevel(varDecl.name)
		if level == 0 && errLevel == nil {
			v.err = fmt.Errorf("variable %s is already declared as parameter", varDecl.name)
			return
		}
	}
	varDecl.expression.accept(v)
	v.varInfo.endVarDecl(varDecl.name)
}

func (v *VariableResolver) visitPrint(printStmt *PrintStatement) {
	printStmt.expression.accept(v)
}

func (v *VariableResolver) visitReturnStmt(returnStmt *ReturnStatement) {
	if !v.inFunctionScope() {
		v.err = fmt.Errorf("return statement is only allowed in function scope")
		return
	}
	if returnStmt.expression != nil {
		returnStmt.expression.accept(v)
	}
}

func (v *VariableResolver) visitExprStmt(exprStmt *ExpressionStatement) {
	exprStmt.expression.accept(v)
}

func (v *VariableResolver) visitIfStmt(ifStmt *IfStatement) {
	ifStmt.condition.accept(v)
	if v.err != nil {
		return
	}
	ifStmt.consequent.accept(v)
	if v.err != nil {
		return
	}
	if ifStmt.alternate != nil {
		ifStmt.alternate.accept(v)
	}
}

func (v *VariableResolver) visitWhileStmt(whileStmt *WhileStatement) {
	whileStmt.condition.accept(v)
	if v.err != nil {
		return
	}
	whileStmt.statement.accept(v)
}

func (v *VariableResolver) visitForStmt(f *ForStatement) {

	v.varInfo = newVarInfo(v.varInfo)
	defer func() {
		v.varInfo = v.varInfo.parent
	}()

	if f.initializer != nil {
		f.initializer.accept(v)
		if v.err != nil {
			return
		}
	}
	if f.condition != nil {
		f.condition.accept(v)
		if v.err != nil {
			return
		}
	}
	if f.increment != nil {
		f.increment.accept(v)
		if v.err != nil {
			return
		}
	}
	f.statement.accept(v)
}

func (v *VariableResolver) visitClassDef(c *ClassDef) {
	v.err = v.varInfo.addName(c.name)
	if v.err != nil {
		return
	}
	v.varInfo = newVarInfo(v.varInfo)
	defer func() {
		v.varInfo = v.varInfo.parent
	}()
	for _, fn := range c.functions {
		fn.accept(v)
		if v.err != nil {
			return
		}
	}
}

func (v *VariableResolver) visitFunctionDef(f *FunctionDef) {
	v.err = v.varInfo.addName(f.name)
	if v.err != nil {
		return
	}
	v.varInfo = newVarInfo(v.varInfo)
	v.varInfo.isParameterInfo = true
	for _, param := range f.parameters {
		v.err = v.varInfo.addName(param)
		if v.err != nil {
			return
		}
	}
	f.body.accept(v)
	v.varInfo = v.varInfo.parent
}

func (v *VariableResolver) visitNumberExpr(*NumberExpr) {}

func (v *VariableResolver) visitBooleanExpr(*BooleanExpr) {}

func (v *VariableResolver) visitNilExpr() {}

func (v *VariableResolver) visitStringExpr(*StringExpr) {}

func (v *VariableResolver) visitIdentifierExpr(identifierExpr *IdentifierExpr) {
	identifierExpr.defLevel, v.err = v.varInfo.getLevel(identifierExpr.name)
}

func (v *VariableResolver) visitGroupExpr(groupExpr *GroupExpr) {
	groupExpr.Inner.accept(v)
}

func (v *VariableResolver) visitUnaryExpr(unaryExpr *UnaryExpr) {
	unaryExpr.Value.accept(v)
}

func (v *VariableResolver) visitBinaryExpr(expr *BinaryExpr) {
	expr.Left.accept(v)
	if v.err != nil {
		return
	}
	expr.Right.accept(v)
}

func (v *VariableResolver) visitAssignment(assignment *Assignment) {
	assignment.right.accept(v)
	if v.err != nil {
		return
	}
	assignment.defLevel, v.err = v.varInfo.getLevel(assignment.left)
}

func (v *VariableResolver) visitCall(call *Call) {
	call.callee.accept(v)
	if v.err != nil {
		return
	}
	for _, arg := range call.args {
		arg.accept(v)
		if v.err != nil {
			return
		}
	}
}

func (v *VariableResolver) inFunctionScope() bool {
	ret := false
	info := v.varInfo
	for {
		if info.isParameterInfo {
			ret = true
			break
		}
		if info.parent == nil {
			break
		}
		info = info.parent
	}
	return ret
}
