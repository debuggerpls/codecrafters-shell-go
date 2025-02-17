package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

const (
	prompt = "$ "
)

func isExec(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func RunShell(stdin io.Reader, stdout io.Writer) int {
	fmt.Fprint(stdout, prompt)

	// Wait for user input
	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		// https://www.gnu.org/software/bash/manual/bash.html#Shell-Operation
		// 1. read input from user's terminal
		input := scanner.Text()

		// 2. break input into words and operators (tokenize), obey quoting rules
		tokens := Tokenize(input)

		// 3. parse tokens into simple and compound commands
		command := tokens[0]
		args := tokens[1:]

		// 4. perform shell expansions

		// 5. perform redirections

		// 6. execute the command
		// 7. optionally wait for command to complete and collect its exit status
		switch command {
		case "exit":
			if len(args) == 0 {
				return 0
			}
			if code, err := strconv.Atoi(args[0]); err == nil {
				return code
			}
			return 0
		case "echo":
			fmt.Fprintln(stdout, strings.Join(args, " "))
		case "pwd":
			pwd, err := os.Getwd()
			if err != nil {
				// TODO: get stderr ?!
				fmt.Fprintf(stdout, "ERROR: %s\n", err)
				return 1
			}
			fmt.Fprintf(stdout, "%s\n", pwd)
		case "cd":
			// TODO: if no args, then change to home
			dest := os.Getenv("HOME")
			if len(args) != 0 && args[0] != "~" {
				dest = args[0]
			}
			err := os.Chdir(dest)
			if err != nil {
				// TODO: stderr, can it be different than no such file or directory?
				fmt.Fprintf(stdout, "cd: %s: No such file or directory\n", args[0])
			}
		case "type":
			argCmd := args[0]
			builtins := []string{"type", "exit", "echo", "pwd", "cd"}

			if slices.Contains(builtins, argCmd) {
				fmt.Fprintf(stdout, "%s is a shell builtin\n", argCmd)
			} else if fPath, err := exec.LookPath(argCmd); err == nil {
				fmt.Fprintf(stdout, "%s is %s\n", argCmd, fPath)
			} else {
				fmt.Fprintf(stdout, "%s: not found\n", argCmd)
			}
		default:
			if isExec(command) {
				cmd := exec.Command(command, args...)
				cmd.Stdout = stdout
				cmd.Stdin = stdin
				if err := cmd.Run(); err != nil {
					fmt.Fprintf(stdout, "ERROR: command run error: %s", err)
				}
			} else {
				fmt.Fprintf(stdout, "%s: command not found\n", command)
			}
		}

		fmt.Fprint(stdout, prompt)
	}

	return 0
}

func Tokenize(input string) []string {
	var tokens []string

	quoted := false
	builder := strings.Builder{}
	const (
		singleQuote = '\''
	)
	lastChar := int32(-1)
	for _, char := range input {
		if char == singleQuote {
			quoted = !quoted
			lastChar = char
			continue
		}
		if !quoted && lastChar == singleQuote && builder.Len() > 0 {
			// a) out of single quotes
			// b) out of single quotes but NOT next to new single quotes
			tokens = append(tokens, builder.String())
			builder.Reset()
		}

		if char == ' ' {
			if lastChar != char {
				// first space
				if quoted {
					builder.WriteRune(char)
				} else {
					if builder.Len() > 0 {
						tokens = append(tokens, builder.String())
						builder.Reset()
					}
				}
				lastChar = char
			} else {
				if quoted {
					builder.WriteRune(char)
				}
			}
			continue
		}

		builder.WriteRune(char)
		lastChar = char
	}

	if builder.Len() > 0 {
		// buffer contains last word
		tokens = append(tokens, builder.String())
	}

	return tokens
}
