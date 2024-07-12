package main

import (
	"fmt"
	"io"
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

	fmt.Println("Token count:", len(tokens))
	fmt.Println("--------")
	for _, t := range tokens {
		fmt.Println(t.ToString())
	}
}
