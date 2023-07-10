package expression_statement

import (
	"MonkeyPL/src/object"
	"MonkeyPL/src/token"
	"fmt"
)

type IdExpression struct {
	token token.Token // 存放变量名
	env   *object.Environment
}

func NewIdExpression(token token.Token, env *object.Environment) (*IdExpression, error) {
	return &IdExpression{token: token, env: env}, nil
}

func (*IdExpression) expressionNode() {}
func (i *IdExpression) Literal() string {
	return i.token.Literal
}
func (i *IdExpression) Name() string {
	return i.token.Literal
}
func (i *IdExpression) String() string {
	return i.token.Literal
}
func (i *IdExpression) Eval() (object.Object, error) {
	if val, ok := i.env.Get(i.token.Literal); !ok {
		return nil, fmt.Errorf("Variable `%s` not found", i.token.Literal)
	} else {
		return val, nil
	}
}
