package expression_statement

import (
	"MonkeyPL/src/common"
	"MonkeyPL/src/object"
	"MonkeyPL/src/token"
)

type InfixExpression interface {
	Expression
	infixExpressionNode()
}
type BaseInfixExpression struct {
	operatorToken token.Token // 运算符词法单元，如+
	left          Expression
	right         Expression
}

func (ie *BaseInfixExpression) expressionNode()      {}
func (ie *BaseInfixExpression) infixExpressionNode() {}
func (ie *BaseInfixExpression) Literal() string {
	return ie.left.Literal() + " " + ie.operatorToken.Literal + " " + ie.right.Literal()
}
func (ie *BaseInfixExpression) String() string {
	return "(" + ie.left.String() + " " + ie.operatorToken.Literal + " " + ie.right.String() + ")"
}
func (ie *BaseInfixExpression) Left() Expression {
	return ie.left
}
func (ie *BaseInfixExpression) Right() Expression {
	return ie.right
}
func (ie *BaseInfixExpression) Operator() token.Token {
	return ie.operatorToken
}
func NewInfixExpression(left Expression, operatorToken token.Token, right Expression) *BaseInfixExpression {
	return &BaseInfixExpression{
		operatorToken: operatorToken,
		left:          left,
		right:         right,
	}
}
func (i *BaseInfixExpression) Eval() (object.Object, error) {
	val, err := i.left.Eval()
	if err != nil {
		return nil, err
	}
	leftVal := val.(object.Integer).Value()
	val, err = i.right.Eval()
	if err != nil {
		return nil, err
	}
	rightVal := val.(object.Integer).Value()

	switch i.operatorToken.Type {
	case token.PLUS:
		return object.NewInteger(leftVal + rightVal), nil
	case token.MINUS:
		return object.NewInteger(leftVal - rightVal), nil
	case token.SLASH:
		return object.NewInteger(leftVal / rightVal), nil
	case token.ASTERISK:
		return object.NewInteger(leftVal * rightVal), nil
	case token.NE:
		return object.NewInteger(common.BoolToInt(leftVal != rightVal)), nil
	case token.EQ:
		return object.NewInteger(common.BoolToInt(leftVal == rightVal)), nil
	case token.GT:
		return object.NewInteger(common.BoolToInt(leftVal > rightVal)), nil
	case token.LT:
		return object.NewInteger(common.BoolToInt(leftVal < rightVal)), nil
	}
	panic("")
}

/*
解析中缀表达式
*/
