package ast

import (
	"MonkeyPL/src/interpreter/config"
	"MonkeyPL/src/interpreter/token"
	"fmt"
)

/*
 */
type LetStatement struct {
	id         *IdExpression
	expression Expression
}

func NewLetStatment(token token.Token, expression Expression) (*LetStatement, error) {
	idExpression, err := NewIdExpression(token, expression.Value())
	if err != nil {
		return nil, fmt.Errorf("Failed when parseing %s", expression)
	}

	return &LetStatement{id: idExpression, expression: expression}, nil
}
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) Literal() string {
	return "let " + ls.id.Literal() + " = " + ls.expression.Literal()
}                                                // 值
func (ls *LetStatement) Expression() Expression  { return ls.expression } // 值
func (ls *LetStatement) Value() config.ValueType { return ls.expression.Value() }
func (ls *LetStatement) Id() *IdExpression       { return ls.id }

/*
 */
type IllegalStatment struct{}

func (is IllegalStatment) statementNode()  {}
func (is IllegalStatment) Literal() string { return "" } // 值
/*
 */
type ReturnStatment struct {
	expression Expression
}

func NewReturnStatment(expression Expression) *ReturnStatment {
	return &ReturnStatment{expression: expression}
}

func (rs *ReturnStatment) statementNode()          {}
func (rs *ReturnStatment) Literal() string         { return "return " + rs.expression.Literal() } // 值
func (rs *ReturnStatment) Expression() Expression  { return rs.expression }                       // 值
func (rs *ReturnStatment) Value() config.ValueType { return rs.expression.Value() }
