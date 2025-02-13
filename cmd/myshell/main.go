package main

import (
	"context"
	"github.com/codecrafters-io/shell-starter-go/pkg/shell"
	"os"
)

func main() {
	code := shell.RunShell(context.Background(), os.Stdin, os.Stdout)
	os.Exit(code)
}
