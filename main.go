package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Welcome to SXL REPL")
		StartRepl(os.Stdin, os.Stdout)
	} else {
		evalFile(os.Args[1])
	}
}

func evalFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to open file %s\n", path)
	}

	l := NewLexer(string(bytes))
	p := NewParser(l)

	program := p.ParseProgram()
	errors := p.Errors()

	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Printf("parser error: %s\n", err)
		}
		return
	}

	env := NewEnvironemnt()
	val := Eval(program, env)
	fmt.Printf("%s\n", val.Inspect())
}
