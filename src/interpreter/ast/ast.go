// ast/ast.go
package ast

import "MonkeyPL/src/interpreter/config"

type Ast interface {
	Literal() string
}

// 表达式
type Statement interface {
	Ast
	statementNode()
}

// 语句
type Expression interface {
	Ast
	expressionNode()
	Value() config.ValueType
}
