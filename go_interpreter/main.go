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
	program, errors := parseFile(path)

	if len(errors) > 0 {
		printParserErrors(errors)
		return
	}

	env := NewEnvironemnt()
	val := Eval(program, env)

	if val.Type() == ERROR_VAL_T {
		fmt.Println(val.Inspect())
	}
}

func parseFile(path string) (program *Program, errors []string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to open file %s\n", path)
	}

	l := NewLexer(string(bytes))
	p := NewParser(l)

	return p.ParseProgram(), p.Errors()
}

func printParserErrors(errors []string) {
	for _, err := range errors {
		fmt.Printf("parser error: %s\n", err)
	}
}

func parseAndPrintFile(path string) {
	program, errors := parseFile(path)

	if len(errors) > 0 {
		printParserErrors(errors)
		return
	}

	fmt.Printf("%s", program.String())
}
