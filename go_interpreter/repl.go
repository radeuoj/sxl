package main

import (
	"bufio"
	"fmt"
	"io"
)

const PROMPT = "> "

func StartRepl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := NewEnvironemnt()

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
			continue
		}

		val := Eval(program, env)
		fmt.Printf("%s\n", val.Inspect())
	}
}
