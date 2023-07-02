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
	//return ok
	//return OK%$
	//return 25

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
		statment, err := p.ParseStatment()
		fmt.Printf("Index %d Literal:  %s\n", i, statment.Literal())
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
		statment, err := p.ParseStatment()
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
		exp, err := p.ParseStatment()
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
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}
