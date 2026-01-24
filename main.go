package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Welcome to MKL REPL")
	StartRepl(os.Stdin, os.Stdout)
}
