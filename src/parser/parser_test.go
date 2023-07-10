package parser

import (
	"MonkeyPL/src/expression_statement"
	"MonkeyPL/src/lexer"
	"fmt"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let z = 21
	$ ssss
	let x = 5
	let y = 10
	let foobar = -838383
	`

	l := lexer.New(input)
	p := New(l)

	tests := []struct {
		expectName    string
		expectedValue int
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
			fmt.Println(statment.String())
		}
	}
}

func TestReturn(t *testing.T) {
	input := `
	return ok
	return OK %$
	return 25
	`

	l := lexer.New(input)
	p := New(l)

	tests := []struct {
		expectedLiteral string
		expectedValue   int
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
			fmt.Println(statment.String())
		}
	}
}

//
//func TestPrefix(t *testing.T) {
//	prefixTests := []string{
//		`let a=!5`,
//		`let, b = 23`,
//		`let c = (a+1)*c+vb`,
//		`let b=-15`,
//	}
//
//	for i, tt := range prefixTests {
//		l := lexer.New(tt)
//		p := New(l)
//		for p.HasNext() {
//			exp, err := p.Next()
//			if err != nil {
//				println(i, " error: ", err.Error())
//			}
//			letStat, _ := exp.(*expression_statement.LetStatement)
//			fmt.Println(letStat)
//		}
//		fmt.Println("===")
//	}
//}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int
		operator   string
		rightValue int
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
		lt, _ := stat.(*expression_statement.LetStatement)
		exp, _ := lt.Expression().(*expression_statement.BaseInfixExpression)
		fmt.Println(exp.String())
	}
}

func TestParsingBoolExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int
		operator   string
		rightValue int
	}{
		{"let a=false+true", 5, "+", 5},
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
		lt, _ := stat.(*expression_statement.LetStatement)
		exp, _ := lt.Expression().(*expression_statement.BaseInfixExpression)
		fmt.Println(exp.String())
	}
}

/*
测试代码变量存放功能
*/
//func TestEnviorment(t *testing.T) {
//	tests := []string{
//		`let a = 5
//let b = a+1`,
//		`let a = 5 * 5
//let b = a*2`,
//		`let a = 5
//let b = a`,
//		`let a = 5
//		let b = a
//		let c = d + b + 5`,
//	}
//	for _, tt := range tests {
//		parser := New(lexer.New(tt))
//		for parser.HasNext() {
//			nxt, err := parser.Next()
//			if err != nil {
//				fmt.Println(err.Error())
//				continue
//			}
//			lt, _ := nxt.(*expression_statement.LetStatement)
//			fmt.Println(lt.Expression().Eval())
//		}
//		fmt.Println()
//	}
//
//}
func TestIfStatement(t *testing.T) {
	tests := []string{
		`if a>b{}`,
		`if a>b{
}else{
let b=1+2
}`,
		`if a<b{
}else if b>a{
}else{
}`,
		`if b>a{
`,
		`if b>a{
}else{}`,
	}
	for i, tt := range tests {
		println(i)
		parser := New(lexer.New(tt))
		nxt, err := parser.Next()
		_, ok := nxt.(expression_statement.NoMoreStatement)
		for !ok {
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			lt, _ := nxt.(*expression_statement.IfStatement)
			nxt, err = parser.Next()
			_, ok = nxt.(expression_statement.NoMoreStatement)
			fmt.Println(lt.Literal())
		}
		fmt.Println()
	}
}
