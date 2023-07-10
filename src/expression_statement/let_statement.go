package expression_statement

import (
	"MonkeyPL/src/object"
	"MonkeyPL/src/token"
	"fmt"
)

type LetStatement struct {
	id         *IdExpression
	expression Expression
}

func NewLetStatment(token token.Token, exp Expression, env *object.Environment) (*LetStatement, error) {
	idExpression, err := NewIdExpression(token, env)
	if err != nil {
		return nil, fmt.Errorf("Failed when parseing %s", exp)
	}

	return &LetStatement{id: idExpression, expression: exp}, nil
}
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) Literal() string {
	return "let " + ls.id.Literal() + " = " + ls.expression.Literal()
}
func (ls *LetStatement) Expression() Expression { return ls.expression }
func (ls *LetStatement) Id() *IdExpression      { return ls.id }
func (ls *LetStatement) String() string {
	return "let " + ls.id.Literal() + " = " + ls.expression.String()
}
