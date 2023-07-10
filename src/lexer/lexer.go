package lexer

import (
	"MonkeyPL/src/token"
	"regexp"
	"strconv"
)

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

type Lexer struct {
	input    string
	pos      int
	curToken token.Token
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.pop()
	for l.expectPeek(token.SPACE) {
		l.pop()
	}
	return l
}

func (l *Lexer) Empty() bool {
	return l.curToken.Type == token.EOF
}

func (l *Lexer) expectPeek(tokenType token.TokenType) bool {
	return tokenType == l.curToken.Type
}

func (l *Lexer) Peek() token.Token {
	return l.curToken
}

func (l *Lexer) Pop() token.Token {
	if l.Empty() {
		return l.curToken
	}
	returnVal := l.curToken
	l.pop()
	for l.expectPeek(token.SPACE) {
		l.pop()
	}

	return returnVal
}

func (l *Lexer) pop() {
	if l.pos == len(l.input) {
		l.curToken = newToken(token.EOF, "")
		return
	}

	res := regs.FindStringSubmatchIndex(l.input[l.pos:])

	if res[0] != 0 {
		l.curToken = newToken(token.ILLEGAL, string(l.input[l.pos]))
		l.pos++
		return
	}

	for i := 2; i < len(res); i += 2 {
		if res[i] == 0 {
			groupName, err := strconv.Atoi(groupNames[i/2])
			if err != nil {
				println(err.Error())
			}
			l.curToken = newToken(token.TokenType(groupName), l.input[l.pos:l.pos+res[i+1]])
			l.pos += res[i+1]
			return
		}
	}
}

func newToken(tokenType token.TokenType, literl string) token.Token {
	return token.Token{Type: tokenType, Literal: literl}
}
