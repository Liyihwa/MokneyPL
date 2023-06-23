package lexer

import (
	"testing"

	"MonkeyPL/src/interpreter/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
   x + y;
};

let result = add(five, ten);
 `

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.ID, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.ID, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.ID, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.ID, "x"},
		{token.COMMA, ","},
		{token.ID, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.ID, "x"},
		{token.PLUS, "+"},
		{token.ID, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.ID, "result"},
		{token.ASSIGN, "="},
		{token.ID, "add"},
		{token.LPAREN, "("},
		{token.ID, "five"},
		{token.COMMA, ","},
		{token.ID, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)
	i := 0
	for l.HasNext() {
		tt := tests[i]
		tok := l.Next()
		if tok.Type == token.SPACE || tok.Type == token.EOL {
			continue
		}
		i++
		if tok.Type != tt.expectedType || tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] -wrong,expected type=%q, got type=%s,expected literal=%s,got literal=%s",
				i, token.Names[tt.expectedType], token.Names[tok.Type], tt.expectedLiteral, tok.Literal)
		}
	}

}
