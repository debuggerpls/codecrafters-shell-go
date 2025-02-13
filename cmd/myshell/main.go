package main

import (
	"github.com/codecrafters-io/shell-starter-go/pkg/shell"
	"os"
)

func main() {
	code := shell.RunShell(os.Getenv, os.Stdin, os.Stdout)
	os.Exit(code)
}
