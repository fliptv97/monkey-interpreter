package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/fliptv97/monkey-interpreter/repl"
)

func main() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello, %s! This is the Monkey programming language!\n", currentUser.Username)
	fmt.Printf("Feel free to type in commands\n")
	err = repl.Start(os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}
}
