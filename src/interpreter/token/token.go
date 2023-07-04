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

var Names = map[TokenType]string{
	ILLEGAL:   "ILLEGAL",
	EOF:       "EOF",
	ID:        "ID",
	LE:        "LE",
	GE:        "GE",
	LT:        "LT",
	GT:        "GT",
	EQ:        "EQ",
	NE:        "NE",
	INT:       "INT",
	ASSIGN:    "ASSIGN",
	PLUS:      "PLUS",
	MINUS:     "MINUS",
	BANG:      "BANG",
	ASTERISK:  "ASTERISK",
	SLASH:     "SLASH",
	BACKSLASH: "BACKSLASH",
	COMMA:     "COMMA",
	SEMICOLON: "SEMICOLON",
	LPAREN:    "LPAREN",
	RPAREN:    "RPAREN",
	LBRACE:    "LBRACE",
	RBRACE:    "RBRACE",
	SPACE:     "SPACE",
	EOL:       "EOL",
	FUNCTION:  "FUNCTION",
	LET:       "LET",
	TRUE:      "TRUE",
	FALSE:     "FALSE",
	IF:        "IF",
	ELSE:      "ELSE",
	RETURN:    "RETURN",
}

var Regs = []struct {
	Type  TokenType
	Regex string
}{
	{INT, `(?:[1-9][0-9]*|0)`},
	{SPACE, `(?:\x20|\t)+`},
	{LE, `<=`},
	{GE, `>=`},
	{EQ, "=="},
	{NE, "!="},
	{LT, `\<`},
	{GT, `\>`},
	{ASSIGN, `=`},
	{PLUS, `\+`},
	{MINUS, `-`},
	{BANG, `!`},
	{ASTERISK, `\*`},
	{SLASH, `/`},
	{BACKSLASH, `\\`},
	{COMMA, `,`},
	{SEMICOLON, `;`},
	{LPAREN, `\(`},
	{RPAREN, `\)`},
	{LBRACE, `\{`},
	{RBRACE, `\}`},
	{EOL, "\n"},
	{FUNCTION, `fn`},
	{LET, `let`},
	{TRUE, `true`},
	{FALSE, `false`},
	{IF, `if`},
	{ELSE, `else`},
	{RETURN, `return`},
	{ID, `[_a-zA-Z][_a-zA-Z0-9]*`},
}

func (t *Token) String() string {
	return fmt.Sprintf("Token{Type: %s;literal: %s}", Names[t.Type], t.Literal)
}
