package repl

import (
	"bufio"
	"fmt"
	"io"
	"mkl/lexer"
	"mkl/parser"
)

const PROMPT = "> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		errors := p.Errors()

		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Printf("parser error: %s\n", err)
			}
		} else {
			fmt.Print(program.String())
		}

		// for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		// 	fmt.Printf("%+v\n", tok)
		// }
	}
}
