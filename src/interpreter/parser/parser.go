package parser

import (
	"MonkeyPL/src/interpreter/ast"
	"MonkeyPL/src/interpreter/lexer"
	"MonkeyPL/src/interpreter/token"
	"errors"
	"fmt"
)

type Parser struct {
	l           *lexer.Lexer
	curStatment ast.Statement
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.parseStatment()
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

// Peek方法返回语法解析器顶
func (p *Parser) Peek() ast.Statement {
	return p.curStatment
}

func (p *Parser) Empty() bool {
	return p.curStatment == nil
}

// parseStatment方法一次解析一整行,将其转化为ast中的语句(Statment),目前只能解析let
func (p *Parser) parseStatment() error {
	for p.expectPeek(token.EOL) {
		p.l.Pop()
	}
	var errorG error
	switch p.l.Peek().Type {
	case token.LET:
		p.curStatment, errorG = p.parseLet()
	case token.EOF:
		p.curStatment = nil
	case token.ILLEGAL: //出现非法token
		p.curStatment = ast.NewIllegalStatment("Unknow token " + p.l.Peek().Literal)
		errorG = errors.New(fmt.Sprintf("%s", p.curStatment))
	}

	return errorG
}

// ParseStatment方法比起parseStatment多了对空行的处理
func (p *Parser) ParseStatment() (ast.Statement, error) {
	temp := p.curStatment
	err := p.parseStatment()
	if err != nil {
		return nil, err
	}
	return temp, nil
}

/* ----------------
*各种实现方法
 */

func (p *Parser) parseLet() (ast.Statement, error) {
	res := ast.LetStatement{Token: p.l.Peek()}
	p.l.Pop() // 跳过let

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
