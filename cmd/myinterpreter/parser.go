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
	var terms []Expr
	var operators []TokenInfo

loop:
	for {
		term, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		terms = append(terms, term)

		token, err := p.peek()
		if err != nil {
			break
		}
		switch token.GetTokenType() {
		case Plus, Minus:
			operators = append(operators, token)
			_, _ = p.advance()
		default:
			break loop
		}
	}

	if len(terms) == 1 {
		return terms[0], nil
	}

	var ret *BinaryExpr

	for i, op := range operators {
		if ret == nil {
			ret = NewBinaryExpr(op, terms[i], terms[i+1])
		} else {
			ret = NewBinaryExpr(op, ret, terms[i+1])
		}
	}

	return ret, nil

}

func (p *Parser) parseTerm() (Expr, error) {
	var factors []Expr
	var operators []TokenInfo

loop:
	for {
		factor, err := p.parseAtomic()
		if err != nil {
			return nil, err
		}
		factors = append(factors, factor)

		token, err := p.peek()
		if err != nil {
			break
		}
		switch token.GetTokenType() {
		case Star, Slash:
			operators = append(operators, token)
			_, _ = p.advance()
		default:
			break loop
		}
	}

	if len(factors) == 1 {
		return factors[0], nil
	}

	var ret *BinaryExpr

	for i, op := range operators {
		if ret == nil {
			ret = NewBinaryExpr(op, factors[i], factors[i+1])
		} else {
			ret = NewBinaryExpr(op, ret, factors[i+1])
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
