package main

import (
	"MonkeyPL/src/lexer"
	"MonkeyPL/src/parser"
	"bufio"
	"io"
	"os"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line) // 一次一行
		p := parser.New(l)
		stat, err := p.Next()
		if err != nil {
			io.WriteString(out, "[ERROR] "+err.Error()+"\n")
			continue
		}

		io.WriteString(out, "> "+stat.String()+"\n")
	}
}

func main() {
	Start(os.Stdin, os.Stdout)
}
