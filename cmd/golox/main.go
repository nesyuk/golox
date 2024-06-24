package main

import (
	"errors"
	"fmt"
	"github.com/nesyuk/golox/lox"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println(errors.New("usage: glox [script]"))
		os.Exit(64)
	} else if len(os.Args) == 2 {
		if err := lox.RunFile(os.Args[1]); err != nil {
			fmt.Printf("failed to read a file: %v\n", err)
		}
	} else {
		lox.RunPrompt()
	}
}
