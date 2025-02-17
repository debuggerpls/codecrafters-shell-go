package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
)

func emptyEnv(key string) string {
	return ""
}

func TestPrompt(t *testing.T) {
	stdout := bytes.Buffer{}
	stdin := bytes.Buffer{}
	RunShell(&stdin, &stdout)

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
	RunShell(stdin, &stdout)

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
	RunShell(stdin, &stdout)

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
			code := RunShell(stdin, &stdout)
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
	RunShell(stdin, &stdout)

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
	stdin := bytes.NewBufferString("type ls\ntype echo\ntype invalid\n")
	RunShell(stdin, &stdout)

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
}

func TestPwd(t *testing.T) {
	t.Run("pwd output", func(t *testing.T) {
		out := bytes.Buffer{}
		in := bytes.NewBufferString("pwd\n")
		exitCode := RunShell(in, &out)
		if exitCode != 0 {
			t.Fatalf("Expected 0 status code")
		}
		got := strings.Trim(out.String(), prompt)

		cmd := exec.Command("pwd")
		wanted := bytes.Buffer{}
		cmd.Stdout = &wanted
		if err := cmd.Run(); err != nil {
			t.Fatalf("os command pwd failed with: %s", err)
		}

		if got != wanted.String() {
			t.Fatalf("got %q wanted %q", got, wanted.String())
		}
	})
	t.Run("pwd as builtin", func(t *testing.T) {
		out := bytes.Buffer{}
		in := bytes.NewBufferString("type pwd\n")
		RunShell(in, &out)
		got := strings.Trim(out.String(), prompt)
		wanted := "pwd is a shell builtin\n"

		if got != wanted {
			t.Fatalf("got %q wanted %q", got, wanted)
		}
	})
}

func TestCd(t *testing.T) {
	t.Run("cd into existing dir", func(t *testing.T) {
		out := bytes.Buffer{}
		in := bytes.NewBufferString("cd /home\n")
		RunShell(in, &out)
		got := strings.Trim(out.String(), prompt)

		wanted := ""
		if got != wanted {
			t.Fatalf("got %q wanted %q", got, wanted)
		}
	})
	t.Run("cd into non-existing", func(t *testing.T) {
		out := bytes.Buffer{}
		in := bytes.NewBufferString("cd /nnnonexistend\n")
		RunShell(in, &out)
		got := strings.Trim(out.String(), prompt)
		wanted := "cd: /nnnonexistend: No such file or directory\n"

		if got != wanted {
			t.Fatalf("got %q wanted %q", got, wanted)
		}
	})
	t.Run("cd to tilde", func(t *testing.T) {
		out := bytes.Buffer{}
		in := bytes.NewBufferString("cd ~\npwd\n")
		RunShell(in, &out)
		got := strings.Trim(out.String(), prompt)
		wanted := os.Getenv("HOME") + "\n"

		if got != wanted {
			t.Fatalf("got %q wanted %q", got, wanted)
		}
	})
}
func TestSingleQuotes(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		tokens []string
	}{
		{"no quotes", "echo shell hello", []string{"echo", "shell", "hello"}},
		{"no quotes multiple spaces", "echo shell  hello", []string{"echo", "shell", "hello"}},
		{"simple quotes", "echo 'shell hello'", []string{"echo", "shell hello"}},
		{"with multiple spaces", "echo 'shell   hello'", []string{"echo", "shell   hello"}},
		{"with multiple spaces and no spaces", "echo 'shell   hello' 'shell''hello'", []string{"echo", "shell   hello", "shellhello"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := Tokenize(testCase.input)
			if !slices.Equal(got, testCase.tokens) {
				t.Errorf("got %q want %q input %q", got, testCase.tokens, testCase.input)
			}
		})
	}
}
