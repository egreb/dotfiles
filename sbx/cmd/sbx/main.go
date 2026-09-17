package main

import (
	"fmt"
	"os"
	"sbx/internal/app"
)

func main() {
	a, e := app.New(os.Args[0], os.Stdin, os.Stdout, os.Stderr)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
	os.Exit(a.Run(os.Args[1:]))
}
