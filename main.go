package main

import (
	"fmt"
	"mkl/repl"
	"os"
)

func main() {
	fmt.Println("Welcome to MKL REPL")
	repl.Start(os.Stdin, os.Stdout)
}
