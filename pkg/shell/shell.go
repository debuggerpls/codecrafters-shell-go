package shell

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
)

func RunShell(ctx context.Context, stdin io.Reader, stdout io.Writer) int {
	fmt.Fprint(stdout, "$ ")

	// Wait for user input
	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if strings.HasPrefix(input, "exit") {
			// TODO: return other values?
			return 0
		}

		fmt.Fprintf(stdout, "%s: command not found\n", scanner.Text())

		fmt.Fprint(stdout, "$ ")
	}

	return 1
}
