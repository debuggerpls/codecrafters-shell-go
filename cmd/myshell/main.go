package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var builtins = map[string]func([]string){
	"exit": exitBuiltin,
}

func exitBuiltin(args []string) {
	// TODO: is it ok to just exit here?
	// TODO: pass exit arg when exiting
	os.Exit(0)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stdout, "ERROR: failed to read: %s\n", err.Error())
			os.Exit(1)
		}

		inputParts := strings.Split(strings.TrimSpace(input), " ")
		command := inputParts[0]

		builtin, ok := builtins[command]
		if ok {
			builtin(inputParts[1:])
		} else {
			fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
		}

	}
}
