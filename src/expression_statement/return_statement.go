package expression_statement

type ReturnStatment struct {
	expression Expression
}

func NewReturnStatment(expression Expression) *ReturnStatment {
	return &ReturnStatment{expression: expression}
}
func (rs *ReturnStatment) statementNode()         {}
func (rs *ReturnStatment) Literal() string        { return "return " + rs.expression.Literal() } // 值
func (rs *ReturnStatment) Expression() Expression { return rs.expression }                       // 值
func (rs *ReturnStatment) String() string         { return "return " + rs.expression.String() }  // 值
