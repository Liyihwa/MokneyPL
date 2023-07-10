package parser

import (
	"MonkeyPL/src/common"
	"MonkeyPL/src/config"
	"MonkeyPL/src/expression_statement"
	"MonkeyPL/src/operator_level"
	"MonkeyPL/src/token"
	"fmt"
	"strconv"
)

func (p *Parser) handleErrorExpression(err error) (*expression_statement.Illegal, error) {
	if err == nil {
		err = fmt.Errorf("Failed when parseing '%s' ", p.l.Peek().Literal)
	}
	for !p.curTokenIs(token.EOL) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}
	p.nextToken()
	return &expression_statement.Illegal{}, err
}

func (p *Parser) parseExpression(precedence config.Precedence) (expression_statement.Expression, error) {
	parseFunc := GetExpressionFunc(p.curToken.Type)
	if parseFunc == nil {
		return p.handleErrorExpression(fmt.Errorf("Failed when parsing expression %s ", p.curToken.Literal))
	}

	leftExp, err := parseFunc(p)
	if err != nil {
		return p.handleErrorExpression(err)
	}
	for !p.isLineEnd() && precedence < p.peekPrecedence() {
		infix := GetInfixExpressionFunc(p.peekToken.Type)
		if infix == nil {
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

func (p *Parser) ParseExpression() (expression_statement.Expression, error) {
	exp, err := p.parseExpression(operator_level.LOWEST)
	p.nextToken()
	return exp, err
}

/*
解析前缀Expression
*/

func parsePrefixExpression(p *Parser) (expression_statement.PrefixExpression, error) {
	prefixToken := p.curToken
	p.nextToken()

	right, err := p.parseExpression(operator_level.PREFIX)
	if err != nil {
		return p.handleErrorExpression(err)
	}
	exp := expression_statement.NewPrefixExpression(prefixToken, right)
	return exp, nil
}

/*
解析中缀表达式
*/
func parseInfixExpression(p *Parser, left expression_statement.Expression) (expression_statement.InfixExpression, error) {
	// 中缀表达式的操作符后面必须有token
	operator := p.curToken
	precedence := p.curPrecedence()
	p.nextToken()

	right, err := p.parseExpression(precedence)
	if err != nil {
		return p.handleErrorExpression(err)
	}
	exp := expression_statement.NewInfixExpression(left, operator, right)
	return exp, nil
}

/*
解析括号内容
*/
func parseGroupedExpression(p *Parser) (expression_statement.Expression, error) {
	p.nextToken() //跳过(
	exp, err := p.parseExpression(operator_level.LOWEST)
	if err != nil {
		return p.handleErrorExpression(err)
	}
	p.nextToken()
	if !p.curTokenIs(token.RPAREN) {
		return p.handleErrorExpression(fmt.Errorf("'( " + exp.Literal() + "' should be followed by ')', not " + p.curToken.Literal))
	}
	return exp, nil
}

/*
解析单个表达式
*/
func parseIdExpression(p *Parser) (expression_statement.Expression, error) {
	return expression_statement.NewIdExpression(p.curToken, &p.env)
}
func parseIntegerExpression(p *Parser) (expression_statement.Expression, error) {
	return expression_statement.NewIntegerExpression(p.curToken)
}

// 因为在monkey中的bool会被替换为 0或1,因此这里干脆返回IntegerExpression
func parseBoolExpression(p *Parser) (expression_statement.Expression, error) {
	return expression_statement.NewIntegerExpression(token.Token{Type: token.INT, Literal: strconv.Itoa(common.BoolToInt(p.curToken.Type == token.TRUE))})
}
