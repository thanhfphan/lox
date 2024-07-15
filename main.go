package main

import (
	"fmt"
	"io"
	"lox/ast"
	"lox/scan"
	"os"
)

func main() {
	file, err := os.Open("main.lox")
	if err != nil {
		panic(err)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	scaner := scan.NewScanner([]rune(string(content)))
	tokens := scaner.ScanTokens()

	parser := ast.NewParser(tokens)
	expr := parser.Parser()
	fmt.Println(expr)
}
