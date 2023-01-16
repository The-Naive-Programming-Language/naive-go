package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"naive/visitor"
	"os"
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
	src, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return run(path, src)
}

func runPrompt() {
	scan := bufio.NewScanner(os.Stdin)
	scan.Split(bufio.ScanLines)

	fmt.Printf("naive> ")
	for scan.Scan() {
		_ = run("<repl>", []byte(scan.Text()))
		fmt.Print("naive> ")
	}
}

func run(filename string, src []byte) error {
	// p := parser.New(token.NewFile(filename), src)

	// p.Parse()

	// for _, stmt := range p.Statements {
	// 	fmt.Println(stmt.String())
	// }

	interp := visitor.NewInterpreter(filename, src)
	interp.Interpret()

	return nil
}
