>   该语言来自书籍[用go语言自制解释器](https://weread.qq.com/web/bookDetail/74d32120813ab6de0g019b0e)
>
>   配套代码 https://interpreterbook.com/waiig_code_1.7.zip

### Monkey语言的特性

```js
let age = 1;
let name = "Monkey";
let result = 10 * (20 / 2);

let thorsten = {"name": "Thorsten", "age": 28};
myArray[0]       // => 1
thorsten["name"] // => "Thorsten"

let add = fn(a, b) { return a + b; };

let fibonacci = fn(x) {
  if (x == 0) {
    0
  } else {
    if (x == 1) {
      1
    } else {
      fibonacci(x - 1) + fibonacci(x - 2);
    }
  }
};

//高阶函数
let twice = fn(f, x) {
  return f(f(x));
};

let addTwo = fn(x) {
  return x + 2;
};
// 我们将addTwo作为变量传入了twice函数中
twice(addTwo, 2); // => 6

```

# 词法分析

>   在进行编译初期，我们需要将源代码分隔为词法单元(token)，然后再构建抽象语法树(AST)。
>
>   
>
>   将源代码分隔为词法单元的过程被称为**词法分析**。词法分析器也叫词法单元生成器(tokenizer)或者扫描器(scanner)

如会被token生成器解析为：

```js
let x = 5 + 5;
```

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

### token

```js
let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
```

如上的代码中，有哪些token？

>   答：let，five，=，5，ten，10，fn，(, x, y, ),{等等都算token

### token分隔

为了灵活和扩展性,我们选择用正则表达式来进行词法分析.

1.   首先我们要定义词法规则:

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

2.   其次,我们定义token变量和Lexer(词法分析器)变量:

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

3.   利用正则表达式进行词法匹配:

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

>   语法分析的作用就是
