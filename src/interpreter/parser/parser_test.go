package parser

import (
	"MonkeyPL/src/interpreter/ast"
	"MonkeyPL/src/interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let z = 21
lol ssss
let x = 5;
let y = 10
let foobar = 838383
`

	l := lexer.New(input)
	p := New(l)

	tests := []struct {
		expectName    string
		expectedValue string
		isError       bool
	}{
		{"z", "21", false},
		{"ssss", "", true},
		{"x", "5", true},
		{"y", "10", false},
		{"foobar", "838383", false},
	}
	for i := 0; !p.Empty(); i++ {
		statment, err := p.ParseStatment()
		if tests[i].isError {
			if err != nil {
				t.Errorf("Line %d should parse failed, but it's not", i)
			}
		} else {
			testLetStatement(t, statment, tests[i].expectName, tests[i].expectedValue)
		}
	}

}

func testLetStatement(t *testing.T, s ast.Statement, name string, value string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Id.Value != value {
		t.Errorf("letStmt.Name.Value should be '%s',not %s", value, letStmt.Id.Value)
		return false
	}

	if letStmt.Id.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() should be '%s',not %s",
			name, letStmt.Id.TokenLiteral())
		return false
	}

	return true
}
