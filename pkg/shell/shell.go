package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

const (
	prompt = "$ "
)

func GetExecutablePath(name string, getenv func(string) string) (string, error) {
	paths := strings.Split(getenv("PATH"), ":")
	for _, p := range paths {
		fPath := filepath.Join(p, name)
		if _, err := os.Stat(fPath); err == nil {
			return fPath, nil
		}
	}
	return "", os.ErrNotExist
}

func RunShell(getenv func(string) string, stdin io.Reader, stdout io.Writer) int {
	fmt.Fprint(stdout, prompt)

	// Wait for user input
	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		// https://www.gnu.org/software/bash/manual/bash.html#Shell-Operation
		// 1. read input from user's terminal
		input := scanner.Text()

		// 2. break input into words and operators (tokenize), obey quoting rules
		tokens := strings.Split(input, " ")

		// 3. parse tokens into simple and compound commands
		command := tokens[0]
		args := tokens[1:]

		// 4. perform shell expansions

		// 5. perform redirections

		// 6. execute the command
		// 7. optionally wait for command to complete and collect its exit status
		switch command {
		case "exit":
			if code, err := strconv.Atoi(args[0]); err == nil {
				return code
			}
			return 0
		case "echo":
			fmt.Fprintln(stdout, strings.Join(args, " "))
		case "type":
			argCmd := args[0]
			builtins := []string{"type", "exit", "echo"}

			if slices.Contains(builtins, argCmd) {
				fmt.Fprintf(stdout, "%s is a shell builtin\n", argCmd)
			} else if fPath, err := GetExecutablePath(argCmd, getenv); err == nil {
				fmt.Fprintf(stdout, "%s is %s\n", argCmd, fPath)
			} else {
				fmt.Fprintf(stdout, "%s: not found\n", argCmd)
			}
		default:
			fmt.Fprintf(stdout, "%s: command not found\n", command)
		}

		fmt.Fprint(stdout, prompt)
	}

	return 1
}
