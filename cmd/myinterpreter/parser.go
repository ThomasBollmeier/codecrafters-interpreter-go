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

func (p *Parser) Parse() (AST, error) {
	return p.parseExpr()
}

func (p *Parser) parseExpr() (Expr, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	var operator TokenInfo
	operator, err = p.peek()
	if err != nil {
		return left, nil
	}

	switch operator.GetTokenType() {
	case Plus, Minus:
		var right Expr
		_, _ = p.advance()
		right, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
		return NewBinaryExpr(operator, left, right), nil
	default:
		return left, nil
	}
}

func (p *Parser) parseTerm() (Expr, error) {
	atomic, err := p.parseAtomic()
	if err != nil {
		return nil, err
	}
	return atomic, nil
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
	default:
		return nil, errors.New(fmt.Sprintf("unexpected token: %s", tt))
	}
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
