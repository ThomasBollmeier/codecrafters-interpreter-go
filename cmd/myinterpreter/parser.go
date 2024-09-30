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
	statements := make([]Statement, 0)

	var stmt Statement

stmts:
	for {
		token, err := p.peek()
		if err != nil {
			return nil, err
		}

		switch token.GetTokenType() {
		case Var:
			stmt, err = p.parseVarDecl()
		case Print:
			stmt, err = p.parsePrintStmt()
		case EOF:
			_, _ = p.advance()
			break stmts
		default:
			stmt, err = p.parseExprStmt()
		}

		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)

	}

	return NewProgram(statements), nil
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
	return p.parseBinary(
		[]TokenType{EqualEqual, BangEqual},
		func() (Expr, error) { return p.parseComparison() })
}

func (p *Parser) parseComparison() (Expr, error) {
	return p.parseBinary(
		[]TokenType{Greater, GreaterEqual, Less, LessEqual},
		func() (Expr, error) { return p.parseSum() })
}

func (p *Parser) parseSum() (Expr, error) {
	return p.parseBinary(
		[]TokenType{Plus, Minus},
		func() (Expr, error) { return p.parseTerm() })
}

func (p *Parser) parseTerm() (Expr, error) {
	return p.parseBinary(
		[]TokenType{Star, Slash},
		func() (Expr, error) { return p.parseAtomic() })
}

func (p *Parser) parseBinary(
	operatorTypes []TokenType,
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

	for i, op := range operators {
		if ret == nil {
			ret = NewBinaryExpr(op, operands[i], operands[i+1])
		} else {
			ret = NewBinaryExpr(op, ret, operands[i+1])
		}
	}

	return ret, nil
}

func (p *Parser) parseAtomic() (Expr, error) {
	token, err := p.advance()
	if err != nil {
		return nil, err
	}
	switch tt := token.GetTokenType(); tt {
	case Number:
		value, _ := strconv.ParseFloat(token.GetLexeme(), 64)
		return NewNumberExpr(value), nil
	case True:
		return NewBooleanExpr(true), nil
	case False:
		return NewBooleanExpr(false), nil
	case Nil:
		return NewNilExpr(), nil
	case String:
		value := strings.Trim(token.GetLexeme(), "\"")
		return NewStringExpr(value), nil
	case Identifier:
		return NewIdentifierExpr(token.GetLexeme()), nil
	case LeftParen:
		return p.parseGroup()
	case Bang, Minus:
		return p.parseUnary(token)
	default:
		return nil, errors.New(fmt.Sprintf("unexpected token: %s", tt))
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
	for _, tokenType := range expected {
		if token.GetTokenType() == tokenType {
			return token, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("unexpected token type &%s", token.GetTokenType()))
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
