package main

import (
	"io"
	"lox/ast"
	"lox/interpret"
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
	interpreter := interpret.New()

	// expr := parser.Parser()
	// interpreter.Interpret(expr)

	stmts := parser.ParserStmt()
	interpreter.InterpretStmt(stmts)

}
