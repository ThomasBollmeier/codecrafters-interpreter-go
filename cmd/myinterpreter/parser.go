package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	done    bool
	scanner *Scanner
	tokens  []TokenInfo
}

func NewParser(content string) *Parser {
	return &Parser{
		done:    false,
		scanner: NewScanner(content),
		tokens:  make([]TokenInfo, 0),
	}
}

func (p *Parser) ParseProgram() (AST, error) {
	statements, err := p.parseDeclarations()
	if err != nil {
		return nil, err
	}

	return NewProgram(statements), nil
}

func (p *Parser) parseDeclarations() ([]Statement, error) {
	statements := make([]Statement, 0)
	var stmt Statement

stmts:
	for {
		token, err := p.peek()
		if err != nil {
			return nil, err
		}

		switch token.GetTokenType() {
		case RightBrace:
			break stmts
		case EOF:
			_, _ = p.advance()
			break stmts
		}

		stmt, err = p.parseDeclaration(token)
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	return statements, nil
}

func (p *Parser) parseDeclaration(nextToken TokenInfo) (Statement, error) {

	var token TokenInfo
	var err error

	if nextToken != nil {
		token = nextToken
	} else {
		token, err = p.peek()
		if err != nil {
			return nil, err
		}
	}

	switch token.GetTokenType() {
	case Var:
		return p.parseVarDecl()
	case Fun:
		return p.parseFunctionDef()
	default:
		return p.parseStatement(token)
	}
}

func (p *Parser) parseStatement(nextToken TokenInfo) (Statement, error) {
	var stmt Statement
	var token TokenInfo
	var err error

	if nextToken != nil {
		token = nextToken
	} else {
		token, err = p.peek()
		if err != nil {
			return nil, err
		}
	}

	switch token.GetTokenType() {
	case Print:
		stmt, err = p.parsePrintStmt()
	case Return:
		stmt, err = p.parseReturnStmt()
	case If:
		stmt, err = p.parseIfStmt()
	case While:
		stmt, err = p.parseWhileStmt()
	case For:
		stmt, err = p.parseForStmt()
	case LeftBrace:
		stmt, err = p.parseBlock()
	default:
		stmt, err = p.parseExprStmt()
	}

	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func (p *Parser) parseReturnStmt() (Statement, error) {
	_, err := p.consume(Return)
	if err != nil {
		return nil, err
	}
	token, err := p.peek()
	if err != nil {
		return nil, err
	}
	if token.GetTokenType() == Semicolon {
		_, _ = p.consume(Semicolon)
		return NewReturnStatement(nil), nil
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(Semicolon)
	if err != nil {
		return nil, err
	}

	return NewReturnStatement(expr), nil
}

func (p *Parser) parseForStmt() (Statement, error) {
	_, err := p.consume(For)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(LeftParen)
	if err != nil {
		return nil, err
	}

	var initializer Statement
	nextToken, err := p.peek()
	if err != nil {
		return nil, err
	}
	if nextToken.GetTokenType() != Semicolon {
		initializer, err = p.parseDeclaration(nil)
		if err != nil {
			return nil, err
		}
		_, isVarDecl := initializer.(*VarDecl)
		_, isExprStmt := initializer.(*ExpressionStatement)
		if !isVarDecl && !isExprStmt {
			return nil, errors.New("for initializer must be var decl. or expression statement")
		}
	} else {
		_, _ = p.advance()
	}

	var condition Expr
	nextToken, err = p.peek()
	if err != nil {
		return nil, err
	}
	if nextToken.GetTokenType() != Semicolon {
		condition, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(Semicolon)
	if err != nil {
		return nil, err
	}

	var increment Expr
	nextToken, err = p.peek()
	if err != nil {
		return nil, err
	}
	if nextToken.GetTokenType() != RightParen {
		increment, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(RightParen)
	if err != nil {
		return nil, err
	}

	statement, err := p.parseStatement(nil)
	if err != nil {
		return nil, err
	}

	return NewForStatement(initializer, condition, increment, statement), nil
}

func (p *Parser) parseWhileStmt() (Statement, error) {
	_, err := p.consume(While)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(LeftParen)
	if err != nil {
		return nil, err
	}
	condition, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(RightParen)
	if err != nil {
		return nil, err
	}

	statement, err := p.parseStatement(nil)
	if err != nil {
		return nil, err
	}

	return NewWhileStatement(condition, statement), nil
}

func (p *Parser) parseIfStmt() (Statement, error) {
	_, err := p.consume(If)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(LeftParen)
	if err != nil {
		return nil, err
	}
	condition, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(RightParen)
	if err != nil {
		return nil, err
	}
	consequent, err := p.parseStatement(nil)
	if err != nil {
		return nil, err
	}
	var alternate Statement = nil
	token, err := p.peek()
	if err == nil && token.GetTokenType() == Else {
		_, _ = p.advance()
		alternate, err = p.parseStatement(nil)
		if err != nil {
			return nil, err
		}
	}

	return NewIfStatement(condition, consequent, alternate), nil
}

func (p *Parser) parseBlock() (AST, error) {
	_, err := p.consume(LeftBrace)
	if err != nil {
		return nil, err
	}

	statements, err := p.parseDeclarations()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RightBrace)
	if err != nil {
		return nil, err
	}

	return NewBlock(statements), nil
}

func (p *Parser) parseFunctionDef() (AST, error) {
	_, err := p.consume(Fun)
	if err != nil {
		return nil, err
	}
	ident, err := p.consume(Identifier)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(LeftParen)
	if err != nil {
		return nil, err
	}
	params, err := p.parseParameters()
	if err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return NewFunctionDef(
		ident.GetLexeme(),
		params,
		*body.(*Block)), nil
}

func (p *Parser) parseParameters() ([]string, error) {
	var parameters []string

	token, err := p.peek()
	if err != nil {
		return nil, err
	}
	if token.GetTokenType() == RightParen {
		_, _ = p.advance()
		return parameters, nil
	}

	for {
		if token.GetTokenType() == Identifier {
			_, _ = p.advance()
			parameters = append(parameters, token.GetLexeme())
		} else {
			return nil, errors.New("expected identifier as parameter")
		}
		token, err = p.peek()
		if err != nil {
			return nil, err
		}
		switch token.GetTokenType() {
		case Comma:
			_, _ = p.advance()
		case RightParen:
			_, _ = p.advance()
			return parameters, nil
		default:
			return nil, errors.New("expected comma or right-paren")
		}
		token, err = p.peek()
		if err != nil {
			return nil, err
		}
	}
}

func (p *Parser) parseVarDecl() (AST, error) {
	_, _ = p.advance() // consume var token
	ident, err := p.consume(Identifier)
	if err != nil {
		return nil, err
	}

	token, err := p.consume(Equal, Semicolon)
	if err != nil {
		return nil, err
	}

	var expr Expr

	if token.GetTokenType() == Equal {

		expr, err = p.parseExpr()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(Semicolon)
		if err != nil {
			return nil, err
		}

	} else {
		expr = NewNilExpr()
	}

	return NewVarDecl(ident.GetLexeme(), expr), nil
}

func (p *Parser) parsePrintStmt() (AST, error) {
	_, _ = p.advance() // consume print token
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	token, err := p.advance()
	if err != nil {
		return nil, err
	}
	if token.GetTokenType() != Semicolon {
		return nil, errors.New("expected Semicolon")
	}

	return NewPrintStatement(expr), nil
}

func (p *Parser) parseExprStmt() (AST, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	token, err := p.advance()
	if err != nil {
		return nil, err
	}
	if token.GetTokenType() != Semicolon {
		return nil, errors.New("expected Semicolon")
	}

	return NewExpressionStatement(expr), nil
}

func (p *Parser) ParseExpression() (AST, error) {
	ast, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	token, err := p.advance()
	if err != nil {
		return nil, err
	}
	if token.GetTokenType() != EOF {
		return nil, errors.New("expected end of tokens")
	}
	return ast, nil
}

func (p *Parser) parseExpr() (Expr, error) {
	nextTokens := p.peekNTokens(2)
	if len(nextTokens) == 2 {
		if nextTokens[0].GetTokenType() == Identifier && nextTokens[1].GetTokenType() == Equal {
			// Assignment:
			_, _ = p.advance()
			_, _ = p.advance()
			ident := nextTokens[0].GetLexeme()
			rhs, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			return NewAssignment(ident, rhs), nil
		} else {
			return p.parseDisjunction()
		}
	} else {
		return p.parseDisjunction()
	}
}

func (p *Parser) parseDisjunction() (Expr, error) {
	return p.parseBinary(
		[]TokenType{Or},
		true,
		func() (Expr, error) { return p.parseConjunction() })
}

func (p *Parser) parseConjunction() (Expr, error) {
	return p.parseBinary(
		[]TokenType{And},
		true,
		func() (Expr, error) { return p.parseEquality() })
}

func (p *Parser) parseEquality() (Expr, error) {
	return p.parseBinary(
		[]TokenType{EqualEqual, BangEqual},
		true,
		func() (Expr, error) { return p.parseComparison() })
}

func (p *Parser) parseComparison() (Expr, error) {
	return p.parseBinary(
		[]TokenType{Greater, GreaterEqual, Less, LessEqual},
		true,
		func() (Expr, error) { return p.parseSum() })
}

func (p *Parser) parseSum() (Expr, error) {
	return p.parseBinary(
		[]TokenType{Plus, Minus},
		true,
		func() (Expr, error) { return p.parseTerm() })
}

func (p *Parser) parseTerm() (Expr, error) {
	return p.parseBinary(
		[]TokenType{Star, Slash},
		true,
		func() (Expr, error) { return p.parseAtomic() })
}

func (p *Parser) parseBinary(
	operatorTypes []TokenType,
	leftAssoc bool,
	parseFunc func() (Expr, error)) (Expr, error) {

	var operands []Expr
	var operators []TokenInfo

	opTypes := map[TokenType]bool{}
	for _, opType := range operatorTypes {
		opTypes[opType] = true
	}

loop:
	for {
		operand, err := parseFunc()
		if err != nil {
			return nil, err
		}
		operands = append(operands, operand)

		token, err := p.peek()
		if err != nil {
			break
		}
		_, isValidOperator := opTypes[token.GetTokenType()]
		if isValidOperator {
			operators = append(operators, token)
			_, _ = p.advance()
		} else {
			break loop
		}
	}

	if len(operands) == 1 {
		return operands[0], nil
	}

	var ret *BinaryExpr

	if leftAssoc {
		for i, op := range operators {
			if ret == nil {
				ret = NewBinaryExpr(op, operands[i], operands[i+1])
			} else {
				ret = NewBinaryExpr(op, ret, operands[i+1])
			}
		}
	} else {
		n := len(operators)
		for i := n - 1; i >= 0; i-- {
			op := operators[i]
			if i == n-1 {
				ret = NewBinaryExpr(op, operands[i], operands[i+1])
			} else {
				ret = NewBinaryExpr(op, operands[i], ret)
			}
		}
	}

	return ret, nil
}

func (p *Parser) parseAtomic() (Expr, error) {
	var expr Expr
	var err error

	token, err := p.advance()
	if err != nil {
		return nil, err
	}
	switch tt := token.GetTokenType(); tt {
	case Number:
		value, _ := strconv.ParseFloat(token.GetLexeme(), 64)
		expr = NewNumberExpr(value)
	case True:
		expr = NewBooleanExpr(true)
	case False:
		expr = NewBooleanExpr(false)
	case Nil:
		expr = NewNilExpr()
	case String:
		value := strings.Trim(token.GetLexeme(), "\"")
		expr = NewStringExpr(value)
	case Identifier:
		expr = NewIdentifierExpr(token.GetLexeme())
	case LeftParen:
		expr, err = p.parseGroup()
	case Bang, Minus:
		expr, err = p.parseUnary(token)
	default:
		return nil, errors.New(fmt.Sprintf("unexpected token: %s", tt))
	}

	if err != nil {
		return nil, err
	}

	nextToken, err := p.peek()
	if err != nil || nextToken.GetTokenType() != LeftParen {
		return expr, nil
	}

	return p.parseCall(expr)
}

func (p *Parser) parseCall(callee Expr) (Expr, error) {
	var args []Expr
	var arg Expr
	var token TokenInfo
	var err error

	_, _ = p.consume(LeftParen)

	token, err = p.peek()
	if err != nil {
		return nil, err
	}
	if token.GetTokenType() == RightParen {
		_, _ = p.advance()
		return &Call{callee, args}, nil
	}

	for {
		arg, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		token, err = p.peek()
		if err != nil {
			return nil, err
		}
		if token.GetTokenType() == RightParen {
			_, _ = p.advance()
			break
		} else {
			_, _ = p.consume(Comma)
		}
	}

	token, err = p.peek()
	if err != nil || token.GetTokenType() != LeftParen {
		return &Call{callee, args}, nil
	} else {
		callee = &Call{callee, args}
		return p.parseCall(callee)
	}
}

func (p *Parser) parseUnary(operator TokenInfo) (Expr, error) {
	value, err := p.parseAtomic()
	if err != nil {
		return nil, err
	}
	return NewUnaryExpr(operator, value), nil
}

func (p *Parser) parseGroup() (Expr, error) {
	inner, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	tok, err := p.advance()
	if err != nil {
		return nil, err
	}
	if tok.GetTokenType() != RightParen {
		return nil, errors.New("expected right paren")
	}
	return NewGroupExpr(inner), nil
}

func (p *Parser) advance() (TokenInfo, error) {
	if p.done {
		return nil, errors.New("no tokens left")
	}

	if len(p.tokens) > 0 {
		var ret = p.tokens[0]
		p.tokens = p.tokens[1:]
		return ret, nil
	}

	return p.scanner.AdvanceToken()
}

func (p *Parser) consume(expected ...TokenType) (TokenInfo, error) {
	token, err := p.advance()
	if err != nil {
		return nil, err
	}

	if len(expected) > 0 {
		for _, tokenType := range expected {
			if token.GetTokenType() == tokenType {
				return token, nil
			}
		}
		return nil, errors.New(fmt.Sprintf("unexpected token type &%s", token.GetTokenType()))
	} else {
		return token, nil
	}
}

func (p *Parser) peek() (TokenInfo, error) {
	if p.done {
		return nil, errors.New("no tokens left")
	}

	if len(p.tokens) == 0 {
		tokenInfo, err := p.scanner.AdvanceToken()
		if err != nil {
			return nil, err
		}
		p.tokens = append(p.tokens, tokenInfo)
	}

	return p.tokens[0], nil
}

func (p *Parser) peekNTokens(n int) []TokenInfo {
	if p.done {
		return nil
	}

	for {
		if len(p.tokens) >= n {
			break
		}
		tokenInfo, err := p.scanner.AdvanceToken()
		if err != nil {
			break
		}
		p.tokens = append(p.tokens, tokenInfo)
	}

	var m int
	if len(p.tokens) >= n {
		m = n
	} else {
		m = len(p.tokens)
	}

	return p.tokens[:m]
}
