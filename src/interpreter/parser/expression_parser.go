package parser

import (
	"MonkeyPL/src/interpreter/ast"
	"MonkeyPL/src/interpreter/config"
	"MonkeyPL/src/interpreter/operator_level"
	"MonkeyPL/src/interpreter/token"
	"fmt"
)

func (p *Parser) handleErrorExpression(err error) (ast.Expression, error) {
	if err == nil {
		err = fmt.Errorf("Failed when parseing `%s` ", p.l.Peek().Literal)
	}
	for !p.curTokenIs(token.EOL) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}
	p.nextToken()
	return &ast.IllegalExpression{}, err
}

func (p *Parser) parseExpression(precedence config.Precedence) (ast.Expression, error) {
	parseFunc := GetExpressionFunc(p.curToken.Type)
	if parseFunc == nil {
		return p.handleErrorExpression(fmt.Errorf("Failed when parsing expression %s", p.curToken.Literal))
	}

	leftExp, err := parseFunc(p)
	if err != nil {
		return p.handleErrorExpression(err)
	}
	for !p.isLineEnd() && precedence < p.peekPrecedence() {
		infix := GetInfixExpressionFunc(p.peekToken.Type)
		if infix == nil { //todo
			return leftExp, nil
		}
		p.nextToken()

		leftExp, err = infix(p, leftExp)
		if err != nil {
			return p.handleErrorExpression(err)
		}
	}

	return leftExp, nil
}

/*
解析前缀Expression
*/

func parsePrefixExpression(p *Parser) (ast.Expression, error) {
	prefixToken := p.curToken
	p.nextToken()

	right, err := p.parseExpression(operator_level.PREFIX)
	if err != nil {
		return nil, err
	}
	expression := ast.NewPrefixExpression(prefixToken, right)
	return expression, nil
}

func parseBangPrefixExpression(p *Parser) (ast.Expression, error) {
	if expression, err := parsePrefixExpression(p); err != nil {
		return p.handleErrorExpression(err)
	} else {
		prefixExpression, _ := expression.(*ast.PrefixExpression)
		return &ast.BangPrefixExpression{PrefixExpression: prefixExpression}, nil
	}
}
func parseMinusPrefixExpression(p *Parser) (ast.Expression, error) {
	if expression, err := parsePrefixExpression(p); err != nil {
		return p.handleErrorExpression(err)
	} else {
		prefixExpression, _ := expression.(*ast.PrefixExpression)
		return &ast.MinusPrefixExpression{PrefixExpression: prefixExpression}, nil
	}
}

/*
解析中缀表达式
*/
// todo
func parseInfixExpression(p *Parser, left ast.Expression) (ast.Expression, error) {
	// 中缀表达式的操作符后面必须有token
	operator := p.curToken
	//if p.isLineEnd() {
	//	return nil, fmt.Errorf("The `%s %s `must followed by an expression", left.Literal(), operator.Literal)
	//}
	precedence := p.curPrecedence()
	p.nextToken()
	right, err := p.parseExpression(precedence)
	if err != nil {
		return p.handleErrorExpression(err)
	}

	expression := ast.NewInfixExpression(left, operator, right)
	return expression, nil
}

/*
解析单个表达式
*/

func parseIdExpression(p *Parser) (ast.Expression, error) {
	return ast.NewIdExpression(p.curToken, 0)
}

func parseIntegerExpression(p *Parser) (ast.Expression, error) {

	return ast.NewIntegerExpression(p.curToken)
}
