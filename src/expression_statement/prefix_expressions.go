package expression_statement

import (
	"MonkeyPL/src/common"
	"MonkeyPL/src/object"
	"MonkeyPL/src/token"
)

type PrefixExpression interface {
	Expression
	prefixExpressionNode()
}
type BasePrefixExpression struct {
	token token.Token
	right Expression
}

func NewPrefixExpression(token token.Token, right Expression) *BasePrefixExpression {

	return &BasePrefixExpression{token: token, right: right}
}
func (*BasePrefixExpression) expressionNode()       {}
func (*BasePrefixExpression) prefixExpressionNode() {}
func (b *BasePrefixExpression) Literal() string {
	return b.token.Literal + b.right.Literal()
}
func (b *BasePrefixExpression) String() string {
	return "(" + b.token.Literal + b.right.Literal() + ")"
}
func (b *BasePrefixExpression) Eval() (object.Object, error) {
	eval, _ := b.right.Eval()
	val := eval.(*object.Integer).Value()
	switch b.token.Type {
	case token.MINUS: //todo

		return object.NewInteger(val), nil
	case token.BANG:
		return object.NewInteger(common.BoolToInt(val == 0)), nil
	}
	panic("")
}
