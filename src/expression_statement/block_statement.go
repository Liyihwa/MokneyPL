package expression_statement

import (
	"MonkeyPL/src/object"
	"bytes"
)

type BlockStatement struct {
	statements []Statement
	env        *object.Environment
}

func NewBlockStatement(statements []Statement, env *object.Environment) BlockStatement {
	return BlockStatement{env: env}
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) Literal() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, s := range bs.statements {
		out.WriteString(s.Literal() + "\n")
	}
	out.WriteString("}")
	return out.String()
}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, s := range bs.statements {
		out.WriteString(s.String() + "\n")
	}
	out.WriteString("}")
	return out.String()
}
