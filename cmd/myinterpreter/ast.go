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

type ReturnStatement struct {
	expression AST
}

func NewReturnStatement(expression AST) *ReturnStatement {
	return &ReturnStatement{expression: expression}
}

func (r *ReturnStatement) accept(visitor AstVisitor) {
	visitor.visitReturnStmt(r)
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

type IfStatement struct {
	condition  Expr
	consequent Statement
	alternate  Statement
}

func NewIfStatement(condition Expr, consequent, alternate Statement) *IfStatement {
	return &IfStatement{condition, consequent, alternate}
}

func (i *IfStatement) accept(visitor AstVisitor) {
	visitor.visitIfStmt(i)
}

type WhileStatement struct {
	condition Expr
	statement Statement
}

func NewWhileStatement(condition Expr, statement Statement) *WhileStatement {
	return &WhileStatement{condition, statement}
}

func (w *WhileStatement) accept(visitor AstVisitor) {
	visitor.visitWhileStmt(w)
}

type ForStatement struct {
	initializer Statement
	condition   Expr
	increment   Expr
	statement   Statement
}

func NewForStatement(initializer Statement, condition Expr, increment, statement Statement) *ForStatement {
	return &ForStatement{initializer, condition, increment, statement}
}

func (f *ForStatement) accept(visitor AstVisitor) {
	visitor.visitForStmt(f)
}

type ClassDef struct {
	name      string
	functions []FunctionDef
}

func NewClassDef(name string) *ClassDef {
	return &ClassDef{name: name}
}

func (c *ClassDef) addFunction(function FunctionDef) {
	c.functions = append(c.functions, function)
}

func (c *ClassDef) accept(visitor AstVisitor) {
	visitor.visitClassDef(c)
}

type FunctionDef struct {
	name       string
	parameters []string
	body       Block
	class      *ClassDef
}

func NewFunctionDef(class *ClassDef, name string, parameters []string, body Block) *FunctionDef {
	return &FunctionDef{name, parameters, body, class}
}

func (f *FunctionDef) accept(visitor AstVisitor) {
	visitor.visitFunctionDef(f)
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
	name     string
	defLevel int // defined <defLevel> levels above the current scope
}

func NewIdentifierExpr(name string) *IdentifierExpr {
	return &IdentifierExpr{name, -1}
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
	left     string
	right    Expr
	defLevel int // LHS defined <defLevel> levels above the current scope
}

func NewAssignment(left string, right Expr) *Assignment {
	return &Assignment{left: left, right: right, defLevel: -1}
}

func (assignment *Assignment) accept(visitor AstVisitor) {
	visitor.visitAssignment(assignment)
}

type Call struct {
	callee Expr
	args   []Expr
}

func NewCall(callee Expr, args []Expr) *Call {
	return &Call{
		callee: callee,
		args:   args,
	}
}

func (call *Call) accept(visitor AstVisitor) {
	visitor.visitCall(call)
}

type AstVisitor interface {
	visitProgram(program *Program)
	visitBlock(block *Block)
	visitVarDecl(varDecl *VarDecl)
	visitPrint(printStmt *PrintStatement)
	visitReturnStmt(returnStmt *ReturnStatement)
	visitExprStmt(exprStmt *ExpressionStatement)
	visitIfStmt(ifStmt *IfStatement)
	visitWhileStmt(whileStmt *WhileStatement)
	visitForStmt(f *ForStatement)
	visitClassDef(c *ClassDef)
	visitFunctionDef(f *FunctionDef)
	visitNumberExpr(numberExpr *NumberExpr)
	visitBooleanExpr(booleanExpr *BooleanExpr)
	visitNilExpr()
	visitStringExpr(stringExpr *StringExpr)
	visitIdentifierExpr(identifierExpr *IdentifierExpr)
	visitGroupExpr(groupExpr *GroupExpr)
	visitUnaryExpr(unaryExpr *UnaryExpr)
	visitBinaryExpr(expr *BinaryExpr)
	visitAssignment(assignment *Assignment)
	visitCall(call *Call)
}
