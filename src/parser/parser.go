package parser

import (
	"MonkeyPL/src/config"
	"MonkeyPL/src/expression_statement"
	"MonkeyPL/src/lexer"
	"MonkeyPL/src/object"
	"MonkeyPL/src/operator_level"
	"MonkeyPL/src/token"
	"fmt"
)

/*
* Parser是语法分析器,他会一次输入一整行代码字符,并返回一个'Statement'
* parseStatment()方法每次返回本行的Statment和error
* Next()方法会调用parseStatment(),并维护curStatment和curError
* Parser在创建(New())时会率先执行parseStatment(),从而将curStatment和curError赋值
*
* 如果在解析某一行出现错误时,那么Parser应该将这一行全部读取完,然后返回 nil的Statment和error
 */

type (
	ExpressionFunc      func(*Parser) (expression_statement.Expression, error)
	InfixExpressionFunc func(*Parser, expression_statement.Expression) (expression_statement.InfixExpression, error)
)

var expressionFuncs = map[token.TokenType]ExpressionFunc{}

var infixExpressionFuncs = map[token.TokenType]InfixExpressionFunc{}

func init() {
	expressionFuncs = make(map[token.TokenType]ExpressionFunc)
	expressionFuncs[token.ID] = parseIdExpression
	expressionFuncs[token.INT] = parseIntegerExpression
	expressionFuncs[token.FALSE] = parseBoolExpression
	expressionFuncs[token.TRUE] = parseBoolExpression

	expressionFuncs[token.BANG] = func(p *Parser) (expression_statement.Expression, error) {
		return parsePrefixExpression(p)
	}
	expressionFuncs[token.MINUS] = func(p *Parser) (expression_statement.Expression, error) {
		return parsePrefixExpression(p)
	}
	expressionFuncs[token.LPAREN] = parseGroupedExpression

	infixExpressionFuncs = make(map[token.TokenType]InfixExpressionFunc)
	infixExpressionFuncs[token.MINUS] = parseInfixExpression
	infixExpressionFuncs[token.PLUS] = parseInfixExpression
	infixExpressionFuncs[token.SLASH] = parseInfixExpression
	infixExpressionFuncs[token.ASTERISK] = parseInfixExpression
	infixExpressionFuncs[token.EQ] = parseInfixExpression
	infixExpressionFuncs[token.NE] = parseInfixExpression
	infixExpressionFuncs[token.LT] = parseInfixExpression
	infixExpressionFuncs[token.GT] = parseInfixExpression

}

func GetExpressionFunc(tokenType token.TokenType) ExpressionFunc {
	return expressionFuncs[tokenType]
}
func GetInfixExpressionFunc(tokenType token.TokenType) InfixExpressionFunc {
	return infixExpressionFuncs[tokenType]
}

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token //在解析Expression用到的中间变量
	peekToken token.Token

	env object.Environment //变量环境
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, env: object.NewEnviorment()}
	p.nextToken()
	p.nextToken()
	for p.curTokenIs(token.EOL) {
		p.nextToken()
	}
	return p
}

func (p *Parser) Next() (expression_statement.Statement, error) {
	if p.l.Empty() {
		return expression_statement.NoMoreStatement{}, fmt.Errorf("No more statments! ")
	}
	return p.parseStatment()
}

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}

func (p *Parser) isLineEnd() bool {
	return p.curToken.Type == token.EOL || p.curToken.Type == token.EOF
}

func (p *Parser) curPrecedence() config.Precedence {
	return operator_level.GetOperatorLevel(p.curToken.Type)
}
func (p *Parser) peekPrecedence() config.Precedence {
	return operator_level.GetOperatorLevel(p.peekToken.Type)
}

func (p *Parser) nextToken() {
	p.curToken, p.peekToken = p.peekToken, p.l.Pop()
}
