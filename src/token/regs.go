package token

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
