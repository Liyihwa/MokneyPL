package parser

import (
	"MonkeyPL/src/interpreter/ast"
	"MonkeyPL/src/interpreter/config"
	"MonkeyPL/src/interpreter/lexer"
	"MonkeyPL/src/interpreter/operator_level"
	"MonkeyPL/src/interpreter/token"
	"fmt"
)

/*
* Parser是语法分析器,他会一次输入一整行代码字符,并返回一个`Statement`
* parseStatment()方法每次返回本行的Statment和error
* Next()方法会调用parseStatment(),并维护curStatment和curError
* Parser在创建(New())时会率先执行parseStatment(),从而将curStatment和curError赋值
*
* 如果在解析某一行出现错误时,那么Parser应该将这一行全部读取完,然后返回 nil的Statment和error
 */

type (
	ExpressionFunc      func(*Parser) (ast.Expression, error)
	InfixExpressionFunc func(*Parser, ast.Expression) (ast.Expression, error)
)

var expressionFuncs = map[token.TokenType]ExpressionFunc{}

var infixExpressionFuncs = map[token.TokenType]InfixExpressionFunc{}

func init() {
	expressionFuncs = make(map[token.TokenType]ExpressionFunc)
	expressionFuncs[token.ID] = parseIdExpression
	expressionFuncs[token.INT] = parseIntegerExpression
	expressionFuncs[token.BANG] = parseBangPrefixExpression
	expressionFuncs[token.MINUS] = parseMinusPrefixExpression

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

	curStatment ast.Statement // 当前语句
	curError    error         // 当前语句的error

	curToken  token.Token //在解析Expression用到的中间变量
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()

	for p.curTokenIs(token.EOL) {
		p.nextToken()
	}

	p.curStatment, p.curError = p.parseStatment()
	return p
}

// Next方法用来维护当前Statment和当前Statment的error
func (p *Parser) Next() (ast.Statement, error) {
	if !p.HasNext() {
		return nil, fmt.Errorf("No more statments! ")
	}
	statment, err := p.curStatment, p.curError
	p.curStatment, p.curError = p.parseStatment()
	return statment, err
}

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}

func (p *Parser) isLineEnd() bool {
	return p.curToken.Type == token.EOL || p.curToken.Type == token.EOF
}

func (p *Parser) HasNext() bool {
	return p.curStatment != nil
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
