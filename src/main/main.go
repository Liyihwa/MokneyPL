package main

import (
	"MonkeyPL/src/interpreter/lexer"
	"MonkeyPL/src/interpreter/token"
)

func main() {
	l := lexer.New("?abc,{}!==")
	for l.HasNext() {
		nxt := l.Next()
		if nxt == nil {
			println("Not Match")
		} else {
			println(token.Names[nxt.Type])
		}
	}
}
