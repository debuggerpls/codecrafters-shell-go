package shell

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestPrompt(t *testing.T) {
	stdout := bytes.Buffer{}
	stdin := bytes.Buffer{}

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	go RunShell(ctx, &stdin, &stdout)

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

func TestREPL(t *testing.T) {
	stdout := bytes.Buffer{}
	stdin := bytes.NewBufferString("invalid_command\ninvalid2\ninvalid3\n")

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	go RunShell(ctx, stdin, &stdout)

	// give some time to start the goroutine
	time.Sleep(time.Millisecond)

	scanner := bufio.NewScanner(&stdout)

	cmds := []string{"invalid_command", "invalid2", "invalid3"}
	for _, cmd := range cmds {
		if !scanner.Scan() {
			t.Fatalf("expected invalid command message")
		}
		want := fmt.Sprintf("$ %s: command not found", cmd)
		got := scanner.Text()
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	}
}
