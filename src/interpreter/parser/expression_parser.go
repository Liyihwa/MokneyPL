package parser

import (
	"MonkeyPL/src/interpreter/ast"
	"MonkeyPL/src/interpreter/operator_level"
	"MonkeyPL/src/interpreter/token"
	"fmt"
)

func (p *Parser) handleErrorExpression(err error) (ast.Expression, error) {
	for !p.expectPeek(token.EOL) && !p.expectPeek(token.EOF) {
		p.l.Pop()
	}
	p.l.Pop()
	return &ast.IllegalExpression{}, err
}

func (p *Parser) parseExpression(precendence int) (ast.Expression, error) {
	prefixFunc := p.prefixParseFunc[p.l.Peek().Type]
	if prefixFunc != nil { //解析前缀Expression,返回PrefixExpression
		return prefixFunc()
	}

	return p.handleErrorExpression(fmt.Errorf("Failed when parseing expression `%s` ", p.l.Peek().Literal))
}

// 前缀Expression
// parser/parser.go

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	prefixToken := p.l.Pop()

	right, err := p.parseExpression(operator_level.PREFIX)
	if err != nil {
		return nil, err
	}
	expression := ast.NewPrefixExpression(prefixToken, right)

	return expression, nil
}

func (p *Parser) parseBangPrefixExpression() (ast.Expression, error) {
	prefixToken := p.l.Pop()

	right, err := p.parseExpression(operator_level.PREFIX)
	if err != nil {
		return nil, err
	}
	expression := &ast.BangPrefixExpression{PrefixExpression: ast.NewPrefixExpression(prefixToken, right)}

	return expression, nil
}

func (p *Parser) parseMinusPrefixExpression() (ast.Expression, error) {
	prefixToken := p.l.Pop()

	right, err := p.parseExpression(operator_level.PREFIX)
	if err != nil {
		return nil, err
	}
	expression := &ast.MinusPrefixExpression{PrefixExpression: ast.NewPrefixExpression(prefixToken, right)}

	return expression, nil
}

func (p *Parser) parseIdExpression() (ast.Expression, error) {
	return ast.NewIdExpression(p.l.Pop(), 0)
}

func (p *Parser) parseIntegerExpression() (ast.Expression, error) {
	return ast.NewIntegerExpression(p.l.Pop())
}
