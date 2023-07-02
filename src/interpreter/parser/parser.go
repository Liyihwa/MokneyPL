package parser

import (
	"MonkeyPL/src/interpreter/ast"
	"MonkeyPL/src/interpreter/lexer"
	"MonkeyPL/src/interpreter/token"
)

/*
* Parser是语法分析器,他会一次输入一整行代码字符,并返回一个`Statement`
* parseStatment()方法每次返回本行的Statment和error
* ParseStatment()方法会调用parseStatment(),并维护curStatment和curError
* Parser在创建(New())时会率先执行parseStatment(),从而将curStatment和curError赋值
*
* 如果在解析某一行出现错误时,那么Parser应该将这一行全部读取完,然后返回 nil的Statment和error
 */

type (
	prefixParseFunc func() (ast.Expression, error)
	infixParseFunc  func(ast.Expression) (ast.Expression, error)
)

type Parser struct {
	l               *lexer.Lexer
	curStatment     ast.Statement // 当前语句
	curError        error         // 当前语句的error
	prefixParseFunc map[token.TokenType]prefixParseFunc
	infixParseFunc  map[token.TokenType]infixParseFunc
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.prefixParseFunc = make(map[token.TokenType]prefixParseFunc)
	p.prefixParseFunc[token.ID] = p.parseIdExpression
	p.prefixParseFunc[token.INT] = p.parseIntegerExpression
	p.prefixParseFunc[token.BANG] = p.parseBangPrefixExpression
	p.prefixParseFunc[token.MINUS] = p.parseMinusPrefixExpression
	p.curStatment, p.curError = p.parseStatment()
	return p
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	return p.l.Peek().Type == tokenType
}

func (p *Parser) expectLineEnd() bool {
	return p.l.Peek().Type == token.EOL || p.l.Peek().Type == token.EOF
}

func (p *Parser) Peek() ast.Statement {
	return p.curStatment
}

func (p *Parser) Empty() bool {
	return p.curStatment == nil
}
