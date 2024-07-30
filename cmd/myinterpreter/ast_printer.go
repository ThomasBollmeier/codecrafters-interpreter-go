package main

import (
	"fmt"
	"math"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) visitNumberExpr(num *NumberExpr) {
	if math.Abs(num.Value-math.Round(num.Value)) > 1e-5 {
		fmt.Printf("%.4f", num.Value)
	} else {
		fmt.Printf("%.1f", num.Value)
	}
}

func (ap *AstPrinter) visitBooleanExpr(be *BooleanExpr) {
	fmt.Printf("%t", be.Value)
}

func (ap *AstPrinter) visitNilExpr() {
	fmt.Printf("nil")
}

func (ap *AstPrinter) visitBinaryExpr(binExpr *BinaryExpr) {
	fmt.Printf("(%s ", binExpr.Operator.GetLexeme())
	binExpr.Left.accept(ap)
	fmt.Printf(" ")
	binExpr.Right.accept(ap)
	fmt.Printf(")\n")
}
