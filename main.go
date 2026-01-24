package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Welcome to SXL REPL")
	StartRepl(os.Stdin, os.Stdout)
}
