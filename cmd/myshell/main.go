package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FIXME: there must be a better solution to overcome initialization loop in typeBuiltin
var builtinsNames = []string{
	"exit",
	"echo",
	"type",
}

var builtins = map[string]func([]string){
	builtinsNames[0]: exitBuiltin,
	builtinsNames[1]: echoBuiltin,
	builtinsNames[2]: typeBuiltin,
}

func IsBuiltin(name string) bool {
	for _, builtin := range builtinsNames {
		if builtin == name {
			return true
		}
	}
	return false
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func IsInPath(name string) (string, bool) {
	path := os.Getenv("PATH")
	for _, p := range strings.Split(path, string(os.PathListSeparator)) {
		fpath := filepath.Join(p, name)
		if IsFile(fpath) {
			return fpath, true
		}
	}

	return "", false
}

func exitBuiltin(args []string) {
	// TODO: is it ok to just exit here?
	// TODO: pass exit arg when exiting
	os.Exit(0)
}

func echoBuiltin(args []string) {
	if len(args) == 0 {
		return
	}

	fmt.Fprintf(os.Stdout, "%s\n", strings.Join(args, " "))
}

func typeBuiltin(args []string) {
	if len(args) == 0 {
		return
	}

	command := args[0]
	if IsBuiltin(command) {
		fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", command)
		return
	}
	if path, ok := IsInPath(command); ok {
		fmt.Fprintf(os.Stdout, "%s is %s\n", command, path)
		return
	}

	fmt.Fprintf(os.Stdout, "%s: not found\n", command)
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
