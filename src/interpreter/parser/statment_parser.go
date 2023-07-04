package parser

import (
	"MonkeyPL/src/interpreter/ast"
	"MonkeyPL/src/interpreter/operator_level"
	"MonkeyPL/src/interpreter/token"

	"fmt"
)

// 出现错误时调用该函数
func (p *Parser) handleErrorStatment(err error) (ast.Statement, error) {
	if err == nil {
		err = fmt.Errorf("Failed when parseing `%s` ", p.l.Peek().Literal)
	}
	for !p.curTokenIs(token.EOL) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}
	p.nextToken()

	return ast.IllegalStatment{}, err
}

func (p *Parser) parseStatment() (ast.Statement, error) {
	for p.curTokenIs(token.EOL) { //跳过空行
		p.nextToken()
	}

	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatment()
	case token.EOF:
		return nil, nil
	case token.RETURN:
		return p.parseReturnStatment()
	case token.ILLEGAL:
		// 出现错误时需要读取一整行
		return p.handleErrorStatment(fmt.Errorf("Unknow token " + p.curToken.Literal))
	}

	return nil, nil
}

/*
	解析一些具体的语句
*/

func (p *Parser) parseLetStatment() (ast.Statement, error) {
	p.nextToken() // 跳过"let"
	var idToken token.Token
	if p.curTokenIs(token.ID) { //解析id
		idToken = p.curToken
	} else {
		return p.handleErrorStatment(fmt.Errorf("The `let` should be followed by a id, not `%s` ", p.curToken.Literal))
	}
	p.nextToken() //跳过"id"

	if !p.curTokenIs(token.ASSIGN) { //判断是否是=
		return p.handleErrorStatment(fmt.Errorf("The `let %s` should be followed by `=`, not `%s` ", idToken.Literal, p.curToken.Literal))
	}
	p.nextToken() //跳过"="
	expression, err := p.parseExpression(operator_level.LOWEST)
	if err != nil {
		return p.handleErrorStatment(err)
	}
	//parseExpression执行后,还需要后移一次token
	p.nextToken()
	letStat, err := ast.NewLetStatment(idToken, expression)

	if err != nil {
		return p.handleErrorStatment(err)
	}

	if !p.isLineEnd() { //判断是否是行尾或者文件末尾
		return p.handleErrorStatment(fmt.Errorf("The `let %s = %s ` should be followed by EOL, not `%s` ", letStat.Id().Name(), letStat.Expression().Literal(), p.curToken.Literal))
	}
	p.nextToken() //跳过 /n

	return letStat, nil
}

func (p *Parser) parseReturnStatment() (ast.Statement, error) {
	p.nextToken() //跳过return

	if expression, err := p.parseExpression(operator_level.LOWEST); err != nil { //解析id
		return p.handleErrorStatment(err)
	} else {
		p.nextToken()
		if !p.isLineEnd() { //判断是否是行尾
			return p.handleErrorStatment(fmt.Errorf("The `return %s` should be followed by EOL(\\n), not `%s` ", expression.Literal(), p.curToken.Literal))
		}
		p.nextToken() //跳过 /n
		return ast.NewReturnStatment(expression), nil
	}
}
