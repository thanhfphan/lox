package main

import (
	"io"
	"lox/ast"
	"lox/interpret"
	"lox/resolve"
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
	// parse statements
	stmts := parser.ParserStmt()
	i := interpret.New()
	// resolving
	resolver := resolve.NewResolver(i)
	resolver.Resolve(stmts)
	i.Interpret(stmts)
}
