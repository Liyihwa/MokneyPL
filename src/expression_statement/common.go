package expression_statement

import (
	"MonkeyPL/src/object"
	"fmt"
)

type Ast interface {
	Literal() string
	String() string // ast的String方法被用于debug
}

// 语句
type Statement interface {
	Ast
	statementNode()
}

type IllegalStatement struct{}

func (is IllegalStatement) statementNode()  {}
func (is IllegalStatement) Literal() string { return "IllegalStatement" } // 值
func (is IllegalStatement) String() string  { return "IllegalStatement" }

type NoMoreStatement struct{}

func (ns NoMoreStatement) statementNode()  {}
func (ns NoMoreStatement) Literal() string { return "NoMoreStatement" } // 值
func (ns NoMoreStatement) String() string  { return "NoMoreStatement" }

// 表达式
type Expression interface {
	Ast
	expressionNode()
	Eval() (object.Object, error)
}

type Illegal struct{}

func (*Illegal) expressionNode()       {}
func (*Illegal) infixExpressionNode()  {}
func (*Illegal) prefixExpressionNode() {}
func (i *Illegal) Literal() string {
	return ""
}
func (i *Illegal) String() string {
	return "ILLEGAL"
}
func (i *Illegal) Eval() (object.Object, error) {
	return nil, fmt.Errorf("It's ILLEGAL EXPRESSION!")
}
