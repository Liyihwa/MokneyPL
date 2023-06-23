package lexer

import (
	"MonkeyPL/src/interpreter/token"
	"regexp"
	"strconv"
)

type Lexer struct {
	input string
	line  int // 当前行数
	pos   int
	regs  *regexp.Regexp
}

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
