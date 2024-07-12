package main

import (
	"errors"
	"fmt"
	"github.com/nesyuk/golox/runtime"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println(errors.New("usage: glox [script]"))
		os.Exit(64)
	} else if len(os.Args) == 2 {
		if err := runtime.RunFile(os.Args[1]); err != nil {
			fmt.Printf("failed to read a file: %v\n", err)
		}
	} else {
		runtime.RunPrompt()
	}
}
