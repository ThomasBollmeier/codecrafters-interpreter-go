package main

type AST interface {
	accept(visitor AstVisitor)
}

type Expr interface {
	AST
}

type NumberExpr struct {
	Value float64
}

func NewNumberExpr(value float64) *NumberExpr {
	return &NumberExpr{value}
}

func (num *NumberExpr) accept(visitor AstVisitor) {
	visitor.visitNumberExpr(num)
}

type BooleanExpr struct {
	Value bool
}

func NewBooleanExpr(value bool) *BooleanExpr {
	return &BooleanExpr{value}
}

func (boolean *BooleanExpr) accept(visitor AstVisitor) {
	visitor.visitBooleanExpr(boolean)
}

type NilExpr struct{}

func NewNilExpr() *NilExpr {
	return &NilExpr{}
}

func (nil *NilExpr) accept(visitor AstVisitor) {
	visitor.visitNilExpr()
}

type StringExpr struct {
	Value string
}

func NewStringExpr(value string) *StringExpr {
	return &StringExpr{value}
}

func (string *StringExpr) accept(visitor AstVisitor) {
	visitor.visitStringExpr(string)
}

type GroupExpr struct {
	Inner Expr
}

func NewGroupExpr(inner Expr) *GroupExpr {
	return &GroupExpr{inner}
}

func (groupExpr *GroupExpr) accept(visitor AstVisitor) {
	visitor.visitGroupExpr(groupExpr)
}

type UnaryExpr struct {
	Operator TokenInfo
	Value    Expr
}

func NewUnaryExpr(operator TokenInfo, value Expr) *UnaryExpr {
	return &UnaryExpr{operator, value}
}

func (unaryExpr *UnaryExpr) accept(visitor AstVisitor) {
	visitor.visitUnaryExpr(unaryExpr)
}

type BinaryExpr struct {
	Left, Right Expr
	Operator    TokenInfo
}

func NewBinaryExpr(operator TokenInfo, left, right Expr) *BinaryExpr {
	return &BinaryExpr{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (binExpr *BinaryExpr) accept(visitor AstVisitor) {
	visitor.visitBinaryExpr(binExpr)
}

type AstVisitor interface {
	visitNumberExpr(numberExpr *NumberExpr)
	visitBooleanExpr(booleanExpr *BooleanExpr)
	visitNilExpr()
	visitStringExpr(stringExpr *StringExpr)
	visitGroupExpr(groupExpr *GroupExpr)
	visitUnaryExpr(unaryExpr *UnaryExpr)
	visitBinaryExpr(expr *BinaryExpr)
}
