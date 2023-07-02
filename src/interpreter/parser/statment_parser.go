package parser

import (
	"MonkeyPL/src/interpreter/ast"
	"MonkeyPL/src/interpreter/operator_level"
	"MonkeyPL/src/interpreter/token"

	"fmt"
)

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
	for !p.expectPeek(token.EOL) && !p.expectPeek(token.EOF) {
		p.l.Pop()
	}
	p.l.Pop()
	return ast.IllegalStatment{}, err
}

func (p *Parser) parseStatment() (ast.Statement, error) {
	for p.expectPeek(token.EOL) { //跳过空行
		p.l.Pop()
	}

	switch p.l.Peek().Type {
	case token.LET:
		return p.parseLetStatment()
	case token.EOF:
		return nil, nil
	case token.RETURN:
		return p.parseReturnStatment()
	case token.ILLEGAL:
		// 出现错误时需要读取一整行
		return p.handleErrorStatment(fmt.Errorf("Unknow token " + p.l.Peek().Literal))
	}

	return nil, nil
}

/*
	解析一些具体的语句
*/

func (p *Parser) parseLetStatment() (ast.Statement, error) {
	p.l.Pop() // 跳过"let"
	var idToken token.Token
	if p.expectPeek(token.ID) { //解析id
		idToken = p.l.Pop()
	} else {
		return p.handleErrorStatment(fmt.Errorf("The `let` should be followed by a id, not `%s` ", p.l.Peek().Literal))
	}

	if !p.expectPeek(token.ASSIGN) { //判断是否是=
		return p.handleErrorStatment(fmt.Errorf("The `let %s` should be followed by `=`, not `%s` ", idToken.Literal, p.l.Peek().Literal))
	}
	p.l.Pop()

	expression, err := p.parseExpression(operator_level.LOWEST)
	if err != nil {
		return p.handleErrorStatment(err)
	}

	letStat, err := ast.NewLetStatment(idToken, expression)
	if err != nil {
		return p.handleErrorStatment(err)
	}

	if !p.expectLineEnd() { //判断是否是行尾或者文件末尾
		return p.handleErrorStatment(fmt.Errorf("The `let %s = %s ` should be followed by EOL, not `%s` ", letStat.Id().Name(), letStat.Expression().Literal(), p.l.Peek().Literal))
	}

	return letStat, nil
}

func (p *Parser) parseReturnStatment() (ast.Statement, error) {
	p.l.Pop() //跳过return

	if expression, err := p.parseExpression(operator_level.LOWEST); err != nil { //解析id
		return nil, err
	} else {
		if !p.expectLineEnd() { //判断是否是行尾
			return p.handleErrorStatment(fmt.Errorf("The `return %s` should be followed by EOL(\\n), not `%s` ", expression.Literal(), p.l.Peek().Literal))
		}
		return ast.NewReturnStatment(expression), nil
	}
}
