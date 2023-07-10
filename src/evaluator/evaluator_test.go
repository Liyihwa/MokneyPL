package evaluator

//func TestEvalIntegerExpression01(t *testing.T) {
//	tests := []struct {
//		input    string
//		expected int64
//	}{
//		{"return a+v*c", 10},
//		{"return 10", 10},
//	}
//
//	for _, tt := range tests {
//		l := lexer.New(tt.input)
//		p := parser.New(l)
//		st, err := p.Next()
//		if err != nil {
//			fmt.Println(err.Error())
//			continue
//		}
//		rt, _ := st.(*expression_statement.ReturnStatment)
//		exp := rt.Expression()
//		fmt.Printf("%s,%T,%d\n", exp.String(), exp.Eval(), exp.Eval())
//	}
//}
//
//func TestEvalIntegerExpression02(t *testing.T) {
//	tests := []struct {
//		input    string
//		expected int64
//	}{
//		{"return 5", 5},
//		{"return 10", 10},
//		{"return -5", -5},
//		{"return -10", -10},
//		{"return 5 + 5 + 5 + 5 - 10", 10},
//		{"return 2 * 2 * 2 * 2 * 2", 32},
//		{"return -50 + 100 + -50", 0},
//		{"return 5 * 2 + 10", 20},
//		{"return 5 + 2 * 10", 25},
//		{"return 20 + 2 * -10", 0},
//		{"return 50 / 2 * 2 + 10", 60},
//		{"return 2 * (5 + 10)", 30},
//		{"return 3 * 3 * 3 + 10", 37},
//		{"return 3 * (3 * 3) + 10", 37},
//		{"return (5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
//	}
//	for _, tt := range tests {
//		l := lexer.New(tt.input)
//		p := parser.New(l)
//		st, err := p.Next()
//		if err != nil {
//			fmt.Println(err.Error())
//			continue
//		}
//		rt, _ := st.(*expression_statement.ReturnStatment)
//		exp := rt.Expression()
//		fmt.Printf("%s,%T,%d\n", exp.String(), exp.Eval(), exp.Eval())
//	}
//}

// evaluator/evaluator_test.go
//func TestEvalBooleanExpression(t *testing.T) {
//	tests := []struct {
//		input    string
//		expected bool
//	}{
//		{"return true", true},
//		{"return false", false},
//		{"return 1 < 2", true},
//		{"return 1 > 2", false},
//		{"return 1 < 1", false},
//		{"return 1 > 1", false},
//		{"return 1 == 1", true},
//		{"return 1 != 1", false},
//		{"return 1 == 2", false},
//		{"return 1 != 2", true},
//	}
//	for _, tt := range tests {
//		l := lexer.New(tt.input)
//		p := parser.New(l)
//		st, err := p.Next()
//		if err != nil {
//			fmt.Println(err.Error())
//			continue
//		}
//		rt, _ := st.(*expression_statement.ReturnStatment)
//		exp := rt.Expression()
//		fmt.Printf("%s,%T,%d\n", exp.String(), exp.Eval(), exp.Eval())
//	}
//}
