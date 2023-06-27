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
				t.Errorf("Sample %d should parse failed, but it's not", i)
			}
		} else {
			testLetStatement(t, statment, tests[i].expectName, tests[i].expectedValue, i)
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string, value string, idx int) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("At index %d, s.TokenLiteral not 'let'. got=%q", idx, s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("At index %d, s not *ast.LetStatement. got=%T", idx, s)
		return false
	}
	//fmt.Println(letStmt.Id.Token)
	//fmt.Println(letStmt.Value)

	if letStmt.Id.Value != value {
		t.Errorf("At index %d, letStmt.Name.Value should be '%s',not %s", idx, value, letStmt.Id.Value)
		return false
	}

	if letStmt.Id.TokenLiteral() != name {
		t.Errorf("At index %d, letStmt.Name.TokenLiteral() should be '%s',not %s",
			idx, name, letStmt.Id.TokenLiteral())
		return false
	}

	return true
}
