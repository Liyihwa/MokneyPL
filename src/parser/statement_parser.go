package parser

import (
	"MonkeyPL/src/expression_statement"
	"MonkeyPL/src/token"

	"fmt"
)

func (p *Parser) skipLine() {

	for !p.l.Empty() && !p.curTokenIs(token.EOL) {
		p.nextToken()
	}
	if !p.l.Empty() {
		p.nextToken()
	}
}

// 出现错误时调用该函数
func (p *Parser) handleErrorStatment(err error) (expression_statement.Statement, error) {
	if err == nil {
		err = fmt.Errorf("Failed when parseing '%s' ", p.l.Peek().Literal)
	}

	return expression_statement.IllegalStatement{}, err
}

func (p *Parser) parseStatment() (expression_statement.Statement, error) {
	for p.curTokenIs(token.EOL) { //跳过空行
		p.nextToken()
	}

	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatment()
	case token.RETURN:
		return p.parseReturnStatment()
	case token.IF:
		return p.parseIfStatment()
	case token.ILLEGAL:
		// 出现错误时需要读取一整行
		return p.handleErrorStatment(fmt.Errorf("Unknow token " + p.curToken.Literal))
	case token.EOF:
		return nil, nil
	}
	return nil, nil
}

/*
	解析一些具体的语句
*/

func (p *Parser) parseLetStatment() (expression_statement.Statement, error) {
	p.nextToken() // 跳过"let"
	var idToken token.Token
	if p.curTokenIs(token.ID) { //解析id
		idToken = p.curToken
	} else {
		return p.handleErrorStatment(fmt.Errorf("The 'let' should be followed by a id, not '%s' ", p.curToken.Literal))
	}
	p.nextToken() //跳过"id"

	if !p.curTokenIs(token.ASSIGN) { //判断是否是=
		return p.handleErrorStatment(fmt.Errorf("The 'let %s' should be followed by '=', not '%s' ", idToken.Literal, p.curToken.Literal))
	}
	p.nextToken() //跳过"="
	expression, err := p.ParseExpression()
	if err != nil {
		return p.handleErrorStatment(err)
	}
	letStat, err := expression_statement.NewLetStatment(idToken, expression, &p.env)

	if err != nil {
		return p.handleErrorStatment(err)
	}

	if !p.isLineEnd() { //判断是否是行尾或者文件末尾
		return p.handleErrorStatment(fmt.Errorf("The 'let %s = %s ' should be followed by EOL, not '%s' ", letStat.Id().Name(), letStat.Expression().Literal(), p.curToken.Literal))
	}
	p.nextToken() //跳过 /n
	eval, err := expression.Eval()
	if err != nil {
		return p.handleErrorStatment(err)
	}
	p.env.Set(idToken.Literal, eval)
	return letStat, nil
}

func (p *Parser) parseReturnStatment() (expression_statement.Statement, error) {
	p.nextToken() //跳过return

	if expression, err := p.ParseExpression(); err != nil { //解析id
		return p.handleErrorStatment(err)
	} else {
		if !p.isLineEnd() { //判断是否是行尾
			return p.handleErrorStatment(fmt.Errorf("The 'return %s' should be followed by EOL(\\n), not '%s' ", expression.Literal(), p.curToken.Literal))
		}
		p.nextToken() //跳过 /n
		return expression_statement.NewReturnStatment(expression), nil
	}
}

func (p *Parser) parseIfStatment() (expression_statement.Statement, error) {
	p.nextToken()                         //跳过 if
	condition, err := p.ParseExpression() //解析 condition
	if err != nil {
		return p.handleErrorStatment(err)
	}

	if !p.curTokenIs(token.LBRACE) {
		return p.handleErrorStatment(fmt.Errorf("'if %s ' should be followed by '{', not %s", condition.Literal(), p.curToken.Literal))
	}
	//解析if{}中的BlockStatement
	tempIfBlock, err := p.parseBlockStatement()
	if err != nil {
		return p.handleErrorStatment(err)
	}

	ifBlock, _ := tempIfBlock.(*expression_statement.BlockStatement)
	if p.curTokenIs(token.ELSE) { //这种情况是if{}else{}或者if{}else if{}
		p.nextToken() //跳过 else

		if p.curTokenIs(token.LBRACE) { //检查是否是else{} 的情况
			tempElseBlock, err := p.parseBlockStatement()
			if err != nil {
				return p.handleErrorStatment(err)
			}
			if !p.isLineEnd() { //解析完 else块后,下一个token应该是EOL
				return p.handleErrorStatment(fmt.Errorf("The 'else' block statement should be followed by EOL, not %s ", p.curToken.Literal))
			}
			p.nextToken() //跳过行尾
			elseBlock, _ := tempElseBlock.(*expression_statement.BlockStatement)
			ifStatment := expression_statement.NewIfStatment(condition, ifBlock, elseBlock, nil)
			return &ifStatment, nil
		} else if p.curTokenIs(token.IF) { //检查是否是 else if{}的情况
			//解析 if表达式
			tempElseIfStatement, err := p.parseIfStatment()
			if err != nil {
				return p.handleErrorStatment(err)
			}

			// 此处无需检查并跳过行尾,因为在` p.parseIfStatment()`中已经检查了
			elseIfStatement, _ := tempElseIfStatement.(*expression_statement.IfStatement)
			ifStatment := expression_statement.NewIfStatment(condition, ifBlock, nil, elseIfStatement)
			return &ifStatment, nil
		} else {
			return p.handleErrorStatment(fmt.Errorf("'else' should be followed by '{' or 'if', not '%s'", p.curToken.Literal))
		}
	}
	//以下是 if{}后不是else的情况
	if !p.isLineEnd() {
		return p.handleErrorStatment(fmt.Errorf("The 'if' block statement should be followed by EOL, not %s ", p.curToken.Literal))
	}
	p.nextToken() //跳过行尾
	ifStatment := expression_statement.NewIfStatment(condition, ifBlock, nil, nil)
	return &ifStatment, nil
}

func (p *Parser) parseBlockStatement() (expression_statement.Statement, error) {
	p.nextToken()                 //跳过 {
	if !p.curTokenIs(token.EOL) { //判断是否是\n
		return p.handleErrorStatment(fmt.Errorf("The '{' should be followed by 'EOL'(\\n), not '%s' ", p.curToken.Literal))
	}
	p.nextToken() // 跳过EOL
	var blockStat []expression_statement.Statement
	for !p.curTokenIs(token.RBRACE) { //只要不是'}'便一直解析
		if p.l.Empty() {
			return p.handleErrorStatment(fmt.Errorf("Not found '}' "))
		}
		fmt.Println(p.curToken)

		stat, err := p.Next()
		if err != nil {
			return p.handleErrorStatment(err)
		}
		blockStat = append(blockStat, stat)
	}
	p.nextToken() //跳过 }
	blockStatement := expression_statement.NewBlockStatement(blockStat, &p.env)
	return &blockStatement, nil
}
