package ast

import (
	"MonkeyPL/src/interpreter/config"
	"MonkeyPL/src/interpreter/token"
	"strconv"
)

/*
 */
type IdExpression struct {
	token token.Token      // 存放变量名
	value config.ValueType // 变量值,目前是int
}

func NewIdExpression(token token.Token, value config.ValueType) (*IdExpression, error) {
	return &IdExpression{token: token, value: value}, nil
}

func (*IdExpression) expressionNode() {}
func (l *IdExpression) Literal() string {
	return l.token.Literal
}
func (l *IdExpression) Name() string {
	return l.token.Literal
}
func (l *IdExpression) Value() config.ValueType {
	return l.value
}

/*
 */
type IntegerExpression struct {
	token token.Token
	value config.ValueType
}

func NewIntegerExpression(token token.Token) (*IntegerExpression, error) {
	value, err := strconv.Atoi(token.Literal)
	if err != nil {
		return nil, err
	}
	return &IntegerExpression{token: token, value: config.ValueType(value)}, nil
}
func (*IntegerExpression) expressionNode() {}
func (l *IntegerExpression) Literal() string {
	return l.token.Literal
}
func (l *IntegerExpression) Value() config.ValueType {
	return l.value
}

/*
 */
type IllegalExpression struct{}

func (i *IllegalExpression) Value() config.ValueType {
	panic("IllegalExpression has no value")
}

func (*IllegalExpression) expressionNode() {}
func (i *IllegalExpression) Literal() string {
	return ""
}

/*
 */
type BangPrefixExpression struct {
	*PrefixExpression
}

func (b *BangPrefixExpression) Value() config.ValueType {
	val := b.PrefixExpression.Value()
	if val == 0 {
		return 1
	} else {
		return 0
	}
}

type MinusPrefixExpression struct {
	*PrefixExpression
}

func (m *MinusPrefixExpression) Value() config.ValueType {
	return -m.PrefixExpression.Value()
}

/*

 */

type PrefixExpression struct {
	token token.Token
	right Expression
}

func NewPrefixExpression(token token.Token, right Expression) *PrefixExpression {
	return &PrefixExpression{token: token, right: right}
}

func (p *PrefixExpression) Value() config.ValueType {
	return p.right.Value()
}

func (*PrefixExpression) expressionNode() {}
func (p *PrefixExpression) Literal() string {
	return p.token.Literal + p.right.Literal()
}
