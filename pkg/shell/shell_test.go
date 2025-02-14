package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func emptyEnv(key string) string {
	return ""
}

func TestPrompt(t *testing.T) {
	stdout := bytes.Buffer{}
	stdin := bytes.Buffer{}
	RunShell(emptyEnv, &stdin, &stdout)

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
	RunShell(emptyEnv, stdin, &stdout)

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
	RunShell(emptyEnv, stdin, &stdout)

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
	testCases := []struct {
		name     string
		stdinStr string
		exitCode int
	}{
		{"exit 0", "exit 0\ncommand\n", 0},
		{"exit non-zero", "exit 1\ncommand\n", 1},
		{"exit without arg", "exit\n", 0},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			stdout := bytes.Buffer{}
			stdin := bytes.NewBufferString(testCase.stdinStr)
			code := RunShell(emptyEnv, stdin, &stdout)
			if code != testCase.exitCode {
				t.Fatalf("expected status code %d, got %d", testCase.exitCode, code)
			}

			scanner := bufio.NewScanner(&stdout)
			if !scanner.Scan() {
				t.Fatalf("expected output, got nothing")
			}
			// expect that REPL breaks after exit
			want := prompt
			got := scanner.Text()
			if got != want {
				t.Errorf("got %q want %q", got, want)
			}
		})
	}
}

func TestEcho(t *testing.T) {
	stdout := bytes.Buffer{}
	stdin := bytes.NewBufferString("echo hello world\n")
	RunShell(emptyEnv, stdin, &stdout)

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
	t.Run("empty path", func(t *testing.T) {
		stdout := bytes.Buffer{}
		stdin := bytes.NewBufferString("type ls\ntype echo\ntype invalid\n")
		RunShell(emptyEnv, stdin, &stdout)

		scanner := bufio.NewScanner(&stdout)
		expected := []string{"$ ls: not found", "$ echo is a shell builtin", "$ invalid: not found"}
		for _, want := range expected {
			if !scanner.Scan() {
				t.Fatalf("expected output")
			}
			got := scanner.Text()
			if got != want {
				t.Fatalf("got %q want %q", got, want)
			}
		}
	})
	t.Run("non-empty path", func(t *testing.T) {
		stdout := bytes.Buffer{}
		stdin := bytes.NewBufferString("type ls\ntype echo\ntype invalid\n")
		env := func(key string) string {
			m := make(map[string]string)
			m["PATH"] = "/usr/bin:/usr/local/bin"

			return m[key]
		}
		RunShell(env, stdin, &stdout)

		scanner := bufio.NewScanner(&stdout)
		expected := []string{"$ ls is /usr/bin/ls", "$ echo is a shell builtin", "$ invalid: not found"}
		for _, want := range expected {
			if !scanner.Scan() {
				t.Fatalf("expected output")
			}
			got := scanner.Text()
			if got != want {
				t.Fatalf("got %q want %q", got, want)
			}
		}
	})
}
