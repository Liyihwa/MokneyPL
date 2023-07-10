package token

import "fmt"

type TokenType int

type Token struct {
	Type    TokenType //类型
	Literal string    //字面量
}

const (
	ILLEGAL = iota
	EOF
	// id
	ID // add, foobar, x, y, ...
	// 整数
	INT // 1343456
	// 运算符
	ASSIGN    //=
	PLUS      //+
	MINUS     // -
	BANG      // !
	ASTERISK  // *
	SLASH     // /
	BACKSLASH // \
	LE
	GE
	LT
	GT
	EQ
	NE
	// 分隔符
	COMMA     // ,
	SEMICOLON // ;
	LPAREN    //(
	RPAREN    //)
	LBRACE    //{
	RBRACE    //}
	SPACE
	EOL
	// 关键字
	FUNCTION // fn
	LET      // let
	TRUE
	FALSE
	IF
	ELSE
	RETURN
)

func (t *Token) String() string {
	return fmt.Sprintf("Token{Type: %s;literal: %s}", Names[t.Type], t.Literal)
}
