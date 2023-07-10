package expression_statement

import (
	"MonkeyPL/src/object"
	"MonkeyPL/src/token"
	"strconv"
)

type IntegerExpression struct {
	Token token.Token
	Val   int
}

func NewIntegerExpression(token token.Token) (*IntegerExpression, error) {
	val, err := strconv.Atoi(token.Literal)
	if err != nil {
		return nil, err
	}
	return &IntegerExpression{Token: token, Val: val}, nil
}

func (*IntegerExpression) expressionNode() {}
func (i *IntegerExpression) Literal() string {
	return i.Token.Literal
}
func (i *IntegerExpression) String() string {
	return i.Literal()
}
func (i *IntegerExpression) Eval() (object.Object, error) {
	return object.NewInteger(i.Val), nil
}
