// ast/ast.go
package ast

import (
	"MonkeyPL/src/interpreter/token"
)

type Ast interface {
	TokenLiteral() string
}

// 表达式
type Statement interface {
	Ast
	statementNode()
}

// 语句
type Expression interface {
	Ast
	expressionNode()
}

/*--------------------------------------------------
* 各种子节点
 */

type LetStatement struct {
	Token token.Token // token.LET词法单元
	Id    *Id
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal } // 值

type IllegalStatment struct{}

func (is IllegalStatment) statementNode()       {}
func (is IllegalStatment) TokenLiteral() string { return "" } // 值

// ------- 定义expression
type Id struct {
	Token token.Token // token.ID词法单元,Token.literal应该是变量名
	Value string      // 变量值
}

func (i *Id) expressionNode()      {}
func (i *Id) TokenLiteral() string { return i.Token.Literal }

func NewId(name string, value string) Id {
	return Id{Token: token.Token{Type: token.ID, Literal: name}, Value: value}
}

type LiteralExpression struct {
	Literal string
}

func (*LiteralExpression) expressionNode() {}
func (l *LiteralExpression) TokenLiteral() string {
	return l.Literal
}
