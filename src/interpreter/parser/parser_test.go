package parser

import (
	"MonkeyPL/src/interpreter/ast"
	"MonkeyPL/src/interpreter/lexer"
	"fmt"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let z = 21
$ ssss
let x = 5;
let y = 10
let foobar = 838383
return ok
return OK%$
return 25
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

	for i := 0; i < len(tests); i++ {
		statment, err := p.ParseStatment()
		if tests[i].isError {
			fmt.Println(err)
			if err == nil {
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

	if letStmt.Id.Value != value {
		t.Errorf("At index %d, letStmt.Name.Value should be '%s',not %s", idx, value, letStmt.Id.Value)
		return false
	}

	if letStmt.Id.TokenLiteral() != name {
		t.Errorf("At index %d, letStmt.Name.TokenLiteral() should be '%s',not %s",
			idx, name, letStmt.Id.TokenLiteral())
		return false
	}
	fmt.Printf("Sample %d passed.\n", idx)
	return true
}

func TestReturn(t *testing.T) {
	input := `
return ok
return OK%$
return 25
`

	l := lexer.New(input)
	p := New(l)

	tests := []struct {
		expectedValue string
		isError       bool
	}{
		{"ok", false},
		{"OK%$", true},
		{"25", false},
	}

	for i := 0; i < len(tests); i++ {
		statment, err := p.ParseStatment()
		if tests[i].isError {
			fmt.Println(err.Error())
			if err == nil {
				t.Errorf("Sample %d should parse failed, but it's not", i)
			}
		} else {
			testReturn(t, statment, tests[i].expectedValue, i)
		}
	}
}

func testReturn(t *testing.T, s ast.Statement, value string, idx int) bool {
	retStatment, ok := s.(*ast.ReturnStatment)
	if !ok {
		t.Errorf("At index %d, s not *ast.ReturnStatment. got=%T", idx, s)
		return false
	}

	if s.TokenLiteral() != "return" {
		t.Errorf("At index %d, s.TokenLiteral not 'return'. got=%q", idx, s.TokenLiteral())
		return false
	}

	if retStatment.Value.TokenLiteral() != value {
		t.Errorf("At index %d, retStatment.Name.Value should be '%s',not %s", idx, value, retStatment.Value.TokenLiteral())
		return false
	}

	fmt.Printf("Sample %d passed.\n", idx)
	return true
}
