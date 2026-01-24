package main

import (
	"bufio"
	"fmt"
	"io"
)

const PROMPT = "> "

func StartRepl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		l := NewLexer(line)
		p := NewParser(l)

		program := p.ParseProgram()
		errors := p.Errors()

		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Printf("parser error: %s\n", err)
			}
		} else {
			fmt.Print(program.String())
		}
	}
}
