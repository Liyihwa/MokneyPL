package parser

import (
	"MonkeyPL/src/interpreter/ast"
	"MonkeyPL/src/interpreter/lexer"
	"MonkeyPL/src/interpreter/token"
	"errors"
	"fmt"
)

/*
* Parser是语法分析器,他会一次输入一整行代码字符,并返回一个`Statement`
* parseStatment()方法每次返回本行的Statment和error
* ParseStatment()方法会调用parseStatment(),并维护curStatment和curError
* Parser在创建(New())时会率先执行ParseStatment(),从而将curStatment和curError赋值
 */

type Parser struct {
	l           *lexer.Lexer
	curStatment ast.Statement // 当前语句
	curError    error         // 当前语句的error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.ParseStatment()
	return p
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	return p.l.Peek().Type == tokenType
}

func (p *Parser) expectPopPeek(tokenType token.TokenType) (token.Token, bool) {
	if p.expectPeek(tokenType) {
		return p.l.Pop(), true
	}
	return p.l.Peek(), false
}

func (p *Parser) Peek() ast.Statement {
	return p.curStatment
}

func (p *Parser) Empty() bool {
	return p.curStatment == nil
}

func (p *Parser) pop() {
	p.l.Pop()
}

// parseStatment方法一次解析一整行,将其转化为ast中的语句(Statment),目前只能解析let
func (p *Parser) parseStatment() (ast.Statement, error) {
	for _, ok := p.expectPopPeek(token.EOL); ok; _, ok = p.expectPopPeek(token.EOL) { //跳过空行

	}
	switch p.l.Peek().Type {
	case token.LET:
		return p.parseLet()
	case token.EOF:
		return nil, nil
	case token.ILLEGAL:
		return ast.NewIllegalStatment("Unknow token " + p.l.Peek().Literal), fmt.Errorf("%s", p.curStatment)
	}
	return nil, nil
}

// ParseStatment方法用来维护当前语句和当前语句的error
func (p *Parser) ParseStatment() (ast.Statement, error) {
	statment, err := p.curStatment, p.curError
	p.curStatment, p.curError = p.parseStatment()
	return statment, err
}

/*
*这里写了不同的语句的实现方法,对于每一种语句(let语句,if语句等等),我们都为其实现一个函数
 */

func (p *Parser) parseLet() (ast.Statement, error) {
	res := ast.LetStatement{Token: p.l.Peek()}
	p.pop() // 跳过let

	idToken := token.Token{}
	if tok, flag := p.expectPopPeek(token.ID); flag { //解析id
		idToken = tok
	} else {
		reason := fmt.Sprintf("The `let` should be followed by a id, not `%s` ", tok.Literal)
		return ast.IllegalStatment{Reason: reason}, errors.New(reason)
	}

	if tok, flag := p.expectPopPeek(token.ASSIGN); !flag { //判断是否是=
		reason := fmt.Sprintf("The `let %s` should be followed by `=`, not `%s` ", idToken.Literal, tok.Literal)
		return ast.IllegalStatment{Reason: reason}, errors.New(reason)
	}

	expression, err := p.parseExpression()
	if err != nil {
		return ast.IllegalStatment{Reason: err.Error()}, err
	}

	temp := ast.NewId(idToken.Literal, expression.TokenLiteral())
	res.Id = &temp
	res.Value = expression

	if tok, flag := p.expectPopPeek(token.EOL); !flag { //判断是否是行尾
		reason := fmt.Sprintf("The `let %s = %s ` should be followed by EOL, not `%s` ", res.Id.Value, res.Value.TokenLiteral(), tok.Literal)
		return ast.IllegalStatment{Reason: reason}, errors.New(reason)
	}

	return &res, nil
}

func (p *Parser) parseExpression() (ast.Expression, error) {
	if tok, flag := p.expectPopPeek(token.INT); flag { //解析字面量
		return &ast.LiteralExpression{Literal: tok.Literal}, nil
	}
	return nil, nil
}
