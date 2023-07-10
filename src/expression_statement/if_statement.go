package expression_statement

import (
	"bytes"
)

type IfStatement struct {
	condition       Expression      //判断if条件的exprssion
	ifBlock         *BlockStatement //if判断成功后要执行的代码块
	elseBlock       *BlockStatement //if判断失败后要执行的代码块,可以为空,可以是ifStatement
	elseIfStatement *IfStatement    //BlockStatement
}

func NewIfStatment(condtion Expression, ifBlock *BlockStatement, elseBlock *BlockStatement, elseStatement *IfStatement) IfStatement {
	return IfStatement{condition: condtion, ifBlock: ifBlock, elseBlock: elseBlock, elseIfStatement: elseStatement}
}
func (is *IfStatement) statementNode() {}
func (is *IfStatement) Literal() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(is.condition.Literal())
	out.WriteString(is.ifBlock.Literal())

	if is.elseBlock != nil {
		out.WriteString(" else")
		out.WriteString(is.elseBlock.Literal())
	}
	if is.elseIfStatement != nil {
		out.WriteString(" else ")
		out.WriteString(is.elseIfStatement.Literal())
	}
	out.WriteString("\n")
	return out.String()
}
func (is *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(is.condition.String())
	out.WriteString(is.ifBlock.String())

	if is.elseBlock != nil {
		out.WriteString(" else")
		out.WriteString(is.elseBlock.String())
	}
	out.WriteString("\n")
	return out.String()
}
