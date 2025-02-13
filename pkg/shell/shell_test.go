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

func TestExit(t *testing.T) {
	stdout := bytes.Buffer{}
	stdin := bytes.NewBufferString("exit 0\ncommand\n")
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	go RunShell(ctx, stdin, &stdout)
	// give some time to start the goroutine
	time.Sleep(time.Millisecond)

	scanner := bufio.NewScanner(&stdout)

	if !scanner.Scan() {
		t.Fatalf("expected output, got nothing")
	}
	// expect that REPL breaks after exit
	want := "$ "
	got := scanner.Text()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestEcho(t *testing.T) {
	stdout := bytes.Buffer{}
	stdin := bytes.NewBufferString("echo hello world\n")
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	go RunShell(ctx, stdin, &stdout)
	// give some time to start the goroutine
	time.Sleep(time.Millisecond)

	scanner := bufio.NewScanner(&stdout)
	if !scanner.Scan() {
		t.Fatalf("expected echo back")
	}
	want := "$ hello world"
	got := scanner.Text()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestType(t *testing.T) {
	stdout := bytes.Buffer{}
	stdin := bytes.NewBufferString("type echo\ntype invalid\n")
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	go RunShell(ctx, stdin, &stdout)
	// give some time to start the goroutine
	time.Sleep(time.Millisecond)

	scanner := bufio.NewScanner(&stdout)
	expected := []string{"$ echo is a shell builtin", "$ invalid: not found"}
	for _, want := range expected {
		if !scanner.Scan() {
			t.Fatalf("expected output")
		}
		got := scanner.Text()
		if got != want {
			t.Fatalf("got %q want %q", got, want)
		}

	}
}
