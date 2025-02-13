package shell

import (
	"bufio"
	"bytes"
	"context"
	"testing"
	"time"
)

func TestPrompt(t *testing.T) {
	stdout := bytes.Buffer{}
	stdin := bytes.NewBufferString("invalid_command\n")

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	go RunShell(ctx, stdin, &stdout)

	// give some time to start the goroutine
	time.Sleep(time.Millisecond)

	scanner := bufio.NewScanner(&stdout)
	if !scanner.Scan() {
		t.Errorf("expected prompt")
	}

	want := "$ "
	got := scanner.Text()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestInvalidCommands(t *testing.T) {
	stdout := bytes.Buffer{}
	stdin := bytes.NewBufferString("invalid_command\n")

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	go RunShell(ctx, stdin, &stdout)

	// give some time to start the goroutine
	time.Sleep(time.Millisecond)

	scanner := bufio.NewScanner(&stdout)
	if !scanner.Scan() {
		t.Fatalf("expected invalid command message")
	}

	want := "$ invalid_command: command not found"
	got := scanner.Text()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
