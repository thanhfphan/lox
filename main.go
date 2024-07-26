package main

import (
	"io"
	"lox/interpreter"
	"lox/parser"
	"lox/resolver"
	"lox/scanner"
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

	tokens := scanner.NewScanner([]rune(string(content))).ScanTokens()
	stmts := parser.New(tokens).ParserStmt()
	i := interpreter.New()
	resolver.NewResolver(i).Resolve(stmts)
	// TODO: check error resolve
	i.Interpret(stmts)
}
