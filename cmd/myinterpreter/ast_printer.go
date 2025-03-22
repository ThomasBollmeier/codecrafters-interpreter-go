package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) visitProgram(*Program) {}

func (ap *AstPrinter) visitBlock(*Block) {}

func (ap *AstPrinter) visitVarDecl(*VarDecl) {}

func (ap *AstPrinter) visitPrint(*PrintStatement) {}

func (ap *AstPrinter) visitReturnStmt(*ReturnStatement) {}

func (ap *AstPrinter) visitExprStmt(*ExpressionStatement) {}

func (ap *AstPrinter) visitIfStmt(*IfStatement) {}

func (ap *AstPrinter) visitWhileStmt(*WhileStatement) {}

func (ap *AstPrinter) visitForStmt(*ForStatement) {}

func (ap *AstPrinter) visitClassDef(*ClassDef) {}

func (ap *AstPrinter) visitFunctionDef(*FunctionDef) {}

func (ap *AstPrinter) visitNumberExpr(num *NumberExpr) {
	numStr := strings.TrimRight(fmt.Sprintf("%f", num.Value), "0")
	if numStr[len(numStr)-1] == uint8('.') {
		numStr = numStr + "0"
	}
	fmt.Print(numStr)
}

func (ap *AstPrinter) visitBooleanExpr(be *BooleanExpr) {
	fmt.Printf("%t", be.Value)
}

func (ap *AstPrinter) visitNilExpr() {
	fmt.Printf("nil")
}

func (ap *AstPrinter) visitStringExpr(str *StringExpr) {
	fmt.Printf("%s", str.Value)
}

func (ap *AstPrinter) visitIdentifierExpr(id *IdentifierExpr) {
	fmt.Printf("id(%s)", id.name)
}

func (ap *AstPrinter) visitGroupExpr(grp *GroupExpr) {
	fmt.Printf("(group ")
	grp.Inner.accept(ap)
	fmt.Printf(")")
}

func (ap *AstPrinter) visitUnaryExpr(unaryExpr *UnaryExpr) {
	opString := unaryExpr.Operator.GetLexeme()
	fmt.Printf("(%s ", opString)
	unaryExpr.Value.accept(ap)
	fmt.Printf(")")
}

func (ap *AstPrinter) visitBinaryExpr(binExpr *BinaryExpr) {
	fmt.Printf("(%s ", binExpr.Operator.GetLexeme())
	binExpr.Left.accept(ap)
	fmt.Printf(" ")
	binExpr.Right.accept(ap)
	fmt.Printf(")")
}

func (ap *AstPrinter) visitAssignment(assignment *Assignment) {
	fmt.Printf("(= %s", assignment.left)
	fmt.Printf(" ")
	assignment.right.accept(ap)
	fmt.Printf(")")
}

func (ap *AstPrinter) visitCall(call *Call) {
	fmt.Printf("(call %s", call.callee)
}
