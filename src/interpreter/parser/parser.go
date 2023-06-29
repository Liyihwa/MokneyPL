package parser

import (
	"MonkeyPL/src/interpreter/ast"
	"MonkeyPL/src/interpreter/lexer"
	"MonkeyPL/src/interpreter/token"
	"fmt"
)

/*
* Parser是语法分析器,他会一次输入一整行代码字符,并返回一个`Statement`
* parseStatment()方法每次返回本行的Statment和error
* ParseStatment()方法会调用parseStatment(),并维护curStatment和curError
* Parser在创建(New())时会率先执行parseStatment(),从而将curStatment和curError赋值
*
* 如果在解析某一行出现错误时,那么Parser应该将这一行全部读取完,然后返回 nil的Statment和error
 */

type Parser struct {
	l           *lexer.Lexer
	curStatment ast.Statement // 当前语句
	curError    error         // 当前语句的error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.curStatment, p.curError = p.parseStatment()
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

// parseStatment方法一次解析一整行,将其转化为ast中的语句(Statment),目前只能解析let
func (p *Parser) parseStatment() (ast.Statement, error) {
	for p.expectPeek(token.EOL) { //跳过空行
		p.l.Pop()
	}

	switch p.l.Peek().Type {
	case token.LET:
		return p.parseLet()
	case token.EOF:
		return nil, nil
	case token.RETURN:
		return p.parseReturn()
	case token.ILLEGAL:
		// 出现错误时需要读取一整行
		return p.handleErrorStatment(fmt.Errorf("Unknow token " + p.l.Peek().Literal))
	}

	return nil, nil
}

// ParseStatment方法用来维护当前Statment和当前Statment的error
func (p *Parser) ParseStatment() (ast.Statement, error) {
	if p.Empty() {
		return nil, fmt.Errorf("No more statments! ")
	}
	statment, err := p.curStatment, p.curError
	p.curStatment, p.curError = p.parseStatment()
	return statment, err
}

// 出现错误时调用该函数
func (p *Parser) handleErrorStatment(err error) (ast.Statement, error) {
	for !p.expectPeek(token.EOL) {
		p.l.Pop()
	}
	p.l.Pop()
	return ast.IllegalStatment{}, err
}
func (p *Parser) handleErrorExpression(err error) (ast.Expression, error) {
	for !p.expectPeek(token.EOL) {
		p.l.Pop()
	}
	p.l.Pop()
	return &ast.IllegalExpression{}, err
}

/*
*这里写了不同的语句的实现方法,对于每一种语句(let语句,if语句等等),我们都为其实现一个函数
 */

func (p *Parser) parseLet() (ast.Statement, error) {
	res := ast.LetStatement{Token: p.l.Peek()}
	p.l.Pop() // 跳过let

	idToken := token.Token{}
	if tok, flag := p.expectPopPeek(token.ID); flag { //解析id
		idToken = tok
	} else {
		return p.handleErrorStatment(fmt.Errorf("The `let` should be followed by a id, not `%s` ", tok.Literal))
	}

	if tok, flag := p.expectPopPeek(token.ASSIGN); !flag { //判断是否是=
		return p.handleErrorStatment(fmt.Errorf("The `let %s` should be followed by `=`, not `%s` ", idToken.Literal, tok.Literal))
	}

	expression, err := p.parseExpression()
	if err != nil {
		return p.handleErrorStatment(err)
	}

	temp := ast.NewId(idToken.Literal, expression.TokenLiteral())
	res.Id = &temp
	res.Value = expression

	if tok, flag := p.expectPopPeek(token.EOL); !flag { //判断是否是行尾
		return p.handleErrorStatment(fmt.Errorf("The `let %s = %s ` should be followed by EOL, not `%s` ", res.Id.TokenLiteral(), res.Value.TokenLiteral(), tok.Literal))
	}

	return &res, nil
}

func (p *Parser) parseReturn() (ast.Statement, error) {
	res := ast.ReturnStatment{Token: p.l.Peek()}

	p.l.Pop()

	if expression, err := p.parseExpression(); err != nil {
		return nil, err
	} else {
		res.Value = expression
	}

	if tok, flag := p.expectPopPeek(token.EOL); !flag { //判断是否是行尾
		return p.handleErrorStatment(fmt.Errorf("The `return %s` should be followed by EOL(\\n), not `%s` ", res.Value.TokenLiteral(), tok.Literal))
	}

	return &res, nil
}

func (p *Parser) parseExpression() (ast.Expression, error) {
	if !(p.expectPeek(token.INT) || p.expectPeek(token.ID)) {
		return p.handleErrorExpression(fmt.Errorf("Failed when parseing expression `%s` ", p.l.Peek().Literal))
	}

	return &ast.LiteralExpression{Literal: p.l.Pop().Literal}, nil
}
