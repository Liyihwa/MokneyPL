package parser

import (
	"MonkeyPL/src/interpreter/ast"
	"MonkeyPL/src/interpreter/config"
	"MonkeyPL/src/interpreter/lexer"
	"fmt"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let z =    21
$ ssss
let x =    5;
let y =    10
let foobar =    -838383
`

	l := lexer.New(input)
	p := New(l)

	tests := []struct {
		expectName    string
		expectedValue config.ValueType
		isError       bool
	}{
		{"z", 21, false},
		{"ssss", 0, true},
		{"x", 5, true},
		{"y", 10, false},
		{"foobar", -838383, false},
	}

	for i := 0; i < len(tests); i++ {
		statment, err := p.Next()
		if tests[i].isError {
			fmt.Printf("At index %d, has error [%s]\n", i, err)
			if err == nil {
				t.Errorf("Sample %d should parse failed, but it's not", i)
			}
		} else {
			testLetStatement(t, statment, tests[i].expectName, tests[i].expectedValue, i)
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string, value config.ValueType, idx int) bool {
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("At index %d, s not *ast.LetStatement. got=%T", idx, s)
		return false
	}

	if letStmt.Value() != value {
		t.Errorf("At index %d, letStmt.Name.Value should be '%d',not %d", idx, value, letStmt.Value())
		return false
	}

	if letStmt.Id().Name() != name {
		t.Errorf("At index %d, letStmt.Name.Literal() should be '%s',not %s",
			idx, name, letStmt.Id().Name())
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
		expectedLiteral string
		expectedValue   config.ValueType
		isError         bool
	}{
		{"ok", 0, false},
		{"OK%$", 0, true},
		{"25", 25, false},
	}

	for i := 0; i < len(tests); i++ {
		statment, err := p.Next()
		if tests[i].isError {
			fmt.Println(err.Error())
			if err == nil {
				t.Errorf("Sample %d should parse failed, but it's not", i)
			}
		} else {
			testReturn(t, statment, tests[i].expectedLiteral, tests[i].expectedValue, i)
		}
	}
}

func testReturn(t *testing.T, s ast.Statement, literal string, value config.ValueType, idx int) bool {
	retStatment, ok := s.(*ast.ReturnStatment)
	if !ok {
		t.Errorf("At index %d, s not *ast.ReturnStatment. got=%T", idx, s)
		return false
	}

	if retStatment.Expression().Literal() != literal {
		t.Errorf("At index %d, retStatment.Expression.Value should be '%s',not %s", idx, literal, retStatment.Expression().Literal())
		return false
	}

	if retStatment.Expression().Value() != value {
		t.Errorf("At index %d, retStatment.Value should be '%d',not %d", idx, value, retStatment.Value())
		return false
	}

	fmt.Printf("Sample %d passed.\n", idx)
	return true
}

// parser/parser_test.go
func TestPrefix(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"let a=!5", "!", 0},
		{"let b=-15", "-", -15},
	}

	for i, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		exp, err := p.Next()
		if err != nil {
			println(i, " error: ", err.Error())
		}
		letStat, _ := exp.(*ast.LetStatement)

		fmt.Printf("transe success: %s %d %T\n", letStat.Expression().Literal(), letStat.Expression().Value(), letStat.Expression())
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  config.ValueType
		operator   string
		rightValue config.ValueType
	}{
		{"let a=a+b+c", 5, "+", 5},
		{"let a=-1 --2", 5, "-", 5},
		{"let a=5 * 5", 5, "*", 5},
		{"let a=5 / 5", 5, "/", 5},
		{"let a=5 > 5", 5, ">", 5},
		{"let a=5 < 5", 5, "<", 5},
		{"let a=5 == 5", 5, "==", 5},
		{"let a=5 <4!= 3>4", 5, "!=", 5},
	}

	for i, tt := range infixTests {
		fmt.Printf("-------------------demo %d\n", i)
		l := lexer.New(tt.input)
		p := New(l)
		stat, err := p.Next()
		if err != nil {
			t.Errorf(err.Error())
		}
		lt, _ := stat.(*ast.LetStatement)
		exp, _ := lt.Expression().(*ast.InfixExpression)
		fmt.Println(exp.String())
	}
}
