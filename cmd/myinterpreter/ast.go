package main

type AST interface {
	accept(visitor AstVisitor)
}

type Statement interface {
	AST
}

type Program struct {
	statements []Statement
}

func NewProgram(statements []Statement) *Program {
	return &Program{statements: statements}
}

func (p *Program) accept(visitor AstVisitor) {
	visitor.visitProgram(p)
}

type Block struct {
	statements []Statement
}

func NewBlock(statements []Statement) *Block {
	return &Block{statements: statements}
}

func (b *Block) accept(visitor AstVisitor) {
	visitor.visitBlock(b)
}

type VarDecl struct {
	name       string
	expression AST
}

func NewVarDecl(name string, expression AST) *VarDecl {
	return &VarDecl{name: name, expression: expression}
}

func (v *VarDecl) accept(visitor AstVisitor) {
	visitor.visitVarDecl(v)
}

type PrintStatement struct {
	expression AST
}

func NewPrintStatement(expression AST) *PrintStatement {
	return &PrintStatement{expression: expression}
}

func (p *PrintStatement) accept(visitor AstVisitor) {
	visitor.visitPrint(p)
}

type ExpressionStatement struct {
	expression AST
}

func NewExpressionStatement(expression AST) *ExpressionStatement {
	return &ExpressionStatement{expression: expression}
}

func (e *ExpressionStatement) accept(visitor AstVisitor) {
	visitor.visitExprStmt(e)
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

type IdentifierExpr struct {
	name string
}

func NewIdentifierExpr(name string) *IdentifierExpr {
	return &IdentifierExpr{name}
}

func (identifier *IdentifierExpr) accept(visitor AstVisitor) {
	visitor.visitIdentifierExpr(identifier)
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

type Assignment struct {
	left  string
	right Expr
}

func NewAssignment(left string, right Expr) *Assignment {
	return &Assignment{left: left, right: right}
}

func (assignment *Assignment) accept(visitor AstVisitor) {
	visitor.visitAssignment(assignment)
}

type AstVisitor interface {
	visitProgram(program *Program)
	visitBlock(block *Block)
	visitVarDecl(varDecl *VarDecl)
	visitPrint(printStmt *PrintStatement)
	visitExprStmt(exprStmt *ExpressionStatement)
	visitNumberExpr(numberExpr *NumberExpr)
	visitBooleanExpr(booleanExpr *BooleanExpr)
	visitNilExpr()
	visitStringExpr(stringExpr *StringExpr)
	visitIdentifierExpr(identifierExpr *IdentifierExpr)
	visitGroupExpr(groupExpr *GroupExpr)
	visitUnaryExpr(unaryExpr *UnaryExpr)
	visitBinaryExpr(expr *BinaryExpr)
	visitAssignment(assignment *Assignment)
}
