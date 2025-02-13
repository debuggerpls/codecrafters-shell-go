package shell

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	prompt = "$ "
)

func RunShell(ctx context.Context, stdin io.Reader, stdout io.Writer) int {
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
		default:
			fmt.Fprintf(stdout, "%s: command not found\n", command)
		}

		fmt.Fprint(stdout, prompt)
	}

	return 1
}
