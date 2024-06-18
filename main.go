package main

import (
	"fmt"
	"os"

	"github.com/beto20/jproject/command"
)

func main() {
	if err := command.Root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
