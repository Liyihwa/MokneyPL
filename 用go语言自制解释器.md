>   该语言来自书籍[用go语言自制解释器](https://weread.qq.com/web/bookDetail/74d32120813ab6de0g019b0e)
>
>   配套代码 https://interpreterbook.com/waiig_code_1.7.zip

# MonkeyPL

```js
let age = 1
let name = "Monkey"
let result = 10 * (20 / 2)

let thorsten = {"name": "Thorsten", "age": 28}
myArray[0]       // => 1
thorsten["name"] // => "Thorsten"

let add = fn(a, b) { return a + b }

let fibonacci = fn(x) {
  if (x == 0) {
    0
  } else {
    if (x == 1) {
      1
    } else {
      fibonacci(x - 1) + fibonacci(x - 2)
    }
  }
}

//高阶函数
let twice = fn(f, x) {
  return f(f(x))
}

let addTwo = fn(x) {
  return x + 2
}
// 我们将addTwo作为变量传入了twice函数中
twice(addTwo, 2) // => 6

```

# 词法分析

>   在进行编译初期，我们需要将源代码分隔为词法单元(token)，然后再构建抽象语法树(AST)。
>
>   
>
>   将源代码分隔为词法单元的过程被称为**词法分析**。词法分析器也叫词法单元生成器(tokenizer)或者扫描器(scanner)

如`let x = 5 + 5;`会被token生成器解析为：

```sql
[
  LET,
  IDENTIFIER("x"),
  EQUAL_SIGN,
  INTEGER(5),
  PLUS_SIGN,
  INTEGER(5),
  SEMICOLON
]
```

## token

一个token就是编程语言中不能再分隔的最小词法单元,例如下边的代码中,包含token:**let，five，=，5，ten，10，fn，(, x, y, ),{**等等

```js
let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
```



一个token应该包含类型以及字面量,因此token的定义如下:

```go
package token

type TokenType int

type Token struct {
	Type    TokenType //类型
	Literal string    //字面量
}
```

此外token也有很多种:

```go
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
```

此外,因为我们选用整数作为TokenType,因此我们还需要一个map从而得到Token名:

```go

var Names = map[TokenType]string{
	ILLEGAL:   "ILLEGAL",
	EOF:       "EOF",
	ID:        "ID",
	LE:        "LE",
	GE:        "GE",
	LT:        "EOF",
	GT:        "EOF",
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
```

### 正则表达式

为了灵活和扩展性,我们选择用正则表达式来进行词法分析.

首先我们要明确有哪些token,以及定义token的词法规则,在Regs中我们定义了可能出现在代码中的token正则表达式

```go
var Regs = []struct {
	Type  TokenType
	Regex string
}{
	{INT, `[\+-]?(?:[1-9][0-9]*|0)`},
	{SPACE, `(?:\x20|\t)+`},
	{LE, `<=`},
	{GE, `>=`},
	{EQ, "=="},
	{NE, "!="},
	{LT, `<`},
	{GT, `>`},
	{ASSIGN, `=`},
	{PLUS, `\+`},
	{MINUS, `-`},
	{BANG, `!`},
	{ASTERISK, `\*`},
	{SLASH, `\\`},
	{BACKSLASH, `/`},
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
```

为了将正则表达式编译,我们采用了包级别的`init()`函数:
在该函数中,我们遍历了Regs,对于每一个正则表达式,我们将其封装在了带有组名的"壳子"中:

```go
var regs *regexp.Regexp
var groupNames []string

func init() {
	var tempRegs []byte

	for _, v := range token.Regs {
		tempRegs = append(tempRegs, []byte("(?P<"+strconv.Itoa(int(v.Type))+">"+v.Regex+")|")...)
	}
	tempRegs = tempRegs[:len(tempRegs)-1]
	regs = regexp.MustCompile(string(tempRegs))
	groupNames = regs.SubexpNames()
}
```

其次,我们定义token变量和Lexer(词法分析器)变量:

```go

type Token struct {
   Type    TokenType //类型
   Literal string    //字面量
   Reg     string    //正则表达式
}


type Lexer struct {
   input string
   line  int // 当前行数
   pos   int
   regs  *regexp.Regexp
}
```

利用正则表达式进行词法匹配:

```go
func New(input string) *Lexer {
	l := &Lexer{input: input}
	var tempRegs []byte
	for _, v := range token.Regs {
		tempRegs = append(tempRegs, []byte("(?P<"+strconv.Itoa(int(v.Type))+">"+v.Regex+")|")...)
	}
	tempRegs = tempRegs[:len(tempRegs)-1]
	l.regs = regexp.MustCompile(string(tempRegs))
	return l
}

func (l *Lexer) HasNext() bool {
	return len(l.input) != l.pos
}

func (l *Lexer) Next() *token.Token {
	res := l.regs.FindStringSubmatchIndex(l.input[l.pos:])
	groupNames := l.regs.SubexpNames()
	if res[0] != 0 {
		l.pos++
		return nil
	}
	for i := 2; i < len(res); i += 2 {
		if res[i] == 0 {
			groupName, err := strconv.Atoi(groupNames[i/2])
			if err != nil {
				println(err.Error())
			}
			tok := newToken(token.TokenType(groupName), l.input[l.pos:l.pos+res[i+1]])
			l.pos += res[i+1]
			return &tok
		}
	}
	return nil
}

func newToken(tokenType token.TokenType, literl string) token.Token {
	return token.Token{Type: tokenType, Literal: literl}
}

```

>   New函数返回了一个词法分析器,在New函数中,我们扫描了`token.Regs`,将每一个token种类都添加到了正则表达式中,最后的正则表达式为:
>
>   `(?P<3>[\+-]?(?:[1-9][0-9]*|0))|(?P<23>(?:\x20|\t)+)|(?P<11><=)|(?P<12>>=)|(?P<15>==)|(?P<16>!=)|(?P<13><)|(?P<14>>)|(?P<4>=)|(?P<5>\+)|(?P<6>-)|(?P<7>!)|(?P<8>\*)|(?P<9>\\)|(?P<10>/)|(?P<17>,)|(?P<18>;)|(?P<19>\()|(?P<20>\))|(?P<21>\{)|(?P<22>\})|(?P<24>\n)|(?P<25>fn)|(?P<26>let)|(?P<27>true)|(?P<28>false)|(?P<29>if)|(?P<30>else)|(?P<31>return)|(?P<2>[_a-zA-Z][_a-zA-Z0-9]*)`

# 语法分析

>   语法分析器将文本或词法单元形式的源代码作为输入，产生一个表示该源代码的数据结构。在建立数据结构时，语法分析器会解析输入，检查其是否符合预期的结构。这个过程就称为语法分析。

我们现在要实现一个语法分析器来将下边的代码转为AST

```js
let x = 5;
let y = 10;
let foobar = add(5, 5);
let barfoo = 5 * 5 / 10 + 18 - add(5, 5) + multiply(124);
let anotherName = barfoo;
```

### 定义表达式和语句

>   在一般的编程语言中,一部分代码可以被分为表达式和语句,表达式是返回一个值的代码片段,而语句则无返回值.
>
>   ```cpp
>   int a=2+3-1;  // 表达式
>   if(a>0){			// 语句
>      ...
>   }
>   ```



```go
// ast/ast.go
package ast

//TokenLiteral()返回该节点对应的字面字符串
type Node interface {
	TokenLiteral() string
}

// 语句
type Statement interface {
	Node
	statementNode()
}

// 表达式
type Expression interface {
	Node
	expressionNode()
}


// 实现语句
type Program struct {
    Statements []Statement
}

func (p *Program) TokenLiteral() string {
    if len(p.Statements) > 0 {
        return p.Statements[0].TokenLiteral()
    } else {
        return ""
    }
}
```

### 定义let语句

```js
let x = 1+2 
```

在这段代码中,有三个部分是需要注意的: let,x还有1+2,因此let语句类需要三个部分:

```go

type LetStatement struct {
	Token token.Token // token.LET词法单元
	Name  *Id
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
```

此外,我们还需要一个表达式来存放变量名:

```go
type Id struct {
	Token token.Token // token.ID词法单元
	Value string
}

func (i *Id) expressionNode()      {}
func (i *Id) TokenLiteral() string { return i.Token.Literal }
```

### 语法分析器



### 解析let

### 解析return

### 解析表达式

自上而下的运算符优先级分析（也称普拉特解析法）
