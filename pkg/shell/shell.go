package shell

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

func RunShell(ctx context.Context, stdin io.Reader, stdout io.Writer) {
	fmt.Fprint(stdout, "$ ")

	// Wait for user input
	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		fmt.Fprintf(stdout, "%s: command not found\n", scanner.Text())

		fmt.Fprint(stdout, "$ ")
	}
}
