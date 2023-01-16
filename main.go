package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"naive/interpreter"
	"naive/parser"
	"naive/token"
)

func main() {
	if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else if len(os.Args) == 1 {
		runPrompt()
	} else {
		printUsage()
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "naive /path/to/script")
	os.Exit(1)
}

func runFile(path string) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	src, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	interp := interpreter.New(path, src)
	interp.Interpret()

	return nil
}

func runPrompt() {
	scan := bufio.NewScanner(os.Stdin)
	scan.Split(bufio.ScanLines)

	interp := interpreter.Default()

	for {
		fmt.Printf("naive> ")
		if !scan.Scan() {
			break
		}
		p := parser.New(token.NewFile("<repl>"), []byte(scan.Text()))
		interp.P = p
		interp.Interpret()
	}
}
