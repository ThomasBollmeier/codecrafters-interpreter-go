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

func (ap *AstPrinter) visitBinaryExpr(binExpr *BinaryExpr) {
	fmt.Printf("(%s ", binExpr.Operator.GetLexeme())
	binExpr.Left.accept(ap)
	fmt.Printf(" ")
	binExpr.Right.accept(ap)
	fmt.Printf(")\n")
}
