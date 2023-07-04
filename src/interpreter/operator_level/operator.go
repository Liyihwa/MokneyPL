package operator_level

import (
	"MonkeyPL/src/interpreter/config"
	"MonkeyPL/src/interpreter/token"
)

/*
定义运算符的优先级
*/

const (
	_ = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var operatorLevelMap = map[token.TokenType]config.Precedence{
	token.EQ:       EQUALS,
	token.NE:       EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

func GetOperatorLevel(tokenType token.TokenType) config.Precedence {
	if level, ok := operatorLevelMap[tokenType]; !ok {
		return LOWEST
	} else {
		return level
	}
}
