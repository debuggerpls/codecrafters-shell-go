package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stdout, "ERROR: failed to read: %s\n", err.Error())
		os.Exit(1)
	}

	command := strings.TrimSpace(input)

	fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
}
