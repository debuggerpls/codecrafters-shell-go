package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	ExitBuiltin = "exit"
	EchoBuiltin = "echo"
	TypeBuiltin = "type"
	PwdBuiltin  = "pwd"
	CdBuiltin   = "cd"
)

var SupportedShellBuiltins = []string{
	ExitBuiltin,
	EchoBuiltin,
	TypeBuiltin,
	PwdBuiltin,
	CdBuiltin,
}

var ShellBuiltins = []Command{
	{
		Name:        ExitBuiltin,
		BuiltinFunc: exitBuiltin,
	},
	{
		Name:        EchoBuiltin,
		BuiltinFunc: echoBuiltin,
	},
	{
		Name:        TypeBuiltin,
		BuiltinFunc: typeBuiltin,
	},
	{
		Name:        PwdBuiltin,
		BuiltinFunc: pwdBuiltin,
	},
	{
		Name:        CdBuiltin,
		BuiltinFunc: cdBuiltin,
	},
}

var NoopCommand = Command{
	BuiltinFunc: func(command *Command) error {
		// do nothing
		return nil
	},
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func GetPath(name string) (string, bool) {
	path := os.Getenv("PATH")
	for _, p := range strings.Split(path, string(os.PathListSeparator)) {
		fpath := filepath.Join(p, name)
		if IsFile(fpath) {
			return fpath, true
		}
	}

	return "", false
}

func exitBuiltin(c *Command) error {
	// TODO: pass exit arg when exiting
	os.Exit(0)
	return nil
}

func echoBuiltin(c *Command) error {
	if len(c.Args) == 0 {
		return nil
	}

	fmt.Fprintf(c.Stdout, "%s\n", strings.Join(c.Args, " "))
	return nil
}

func typeBuiltin(c *Command) error {
	if len(c.Args) == 0 {
		return nil
	}

	command := c.Args[0]
	for _, builtin := range SupportedShellBuiltins {
		if command == builtin {
			fmt.Fprintf(c.Stdout, "%s is a shell builtin\n", command)
			return nil
		}
	}
	if path, ok := GetPath(command); ok {
		fmt.Fprintf(c.Stdout, "%s is %s\n", command, path)
		return nil
	}

	return fmt.Errorf("%s: not found", command)
}

func pwdBuiltin(c *Command) error {
	//fmt.Fprintf(os.Stdout, "%s\n", os.Getenv("PWD"))
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Fprintf(c.Stdout, "%s\n", pwd)
	return nil
}

func cdBuiltin(c *Command) error {
	var dir string
	if len(c.Args) == 0 || c.Args[0] == "~" {
		dir = os.Getenv("HOME")
	} else {
		dir = c.Args[0]
	}

	if err := os.Chdir(dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("cd: %s: No such file or directory", dir)
		} else {
			// TODO: do not pass the underlying error
			return fmt.Errorf("cd: %s: %s", dir, err)
		}
	}

	return nil
}

type Command struct {
	Name        string
	Path        string
	Args        []string
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
	BuiltinFunc func(*Command) error
}

func ParseInput(in string) (*Command, error) {
	var command *Command

	parts := strings.Fields(in)

	if len(parts) == 0 {
		return &NoopCommand, nil
	}

	commandName := parts[0]

	for _, builtin := range ShellBuiltins {
		if commandName == builtin.Name {
			command = &builtin
			break
		}
	}

	if command == nil {
		commandPath, ok := GetPath(commandName)
		if ok {
			command = &Command{
				Name: commandName,
				Path: commandPath,
			}
		}
	}

	if command == nil {
		return nil, fmt.Errorf("%s: command not found", commandName)
	}

	quoteCount := strings.Count(in, "'") + strings.Count(in, "\"")
	if quoteCount == 0 {
		command.Args = parts[1:]
	} else {
		// FIXME: check for invalid quoteCount, like not even number
		// FIXME: cases like CMD'ASDASD' or CMD'ASda ' etc
		argString := strings.TrimSpace(in[strings.Index(in, " "):])

		var prevC rune
		var prevQuote rune
		inQuotes := false
		quoteIndex := 0
		for i, c := range argString {
			//fmt.Printf("Index=%d, char=%c (%d)\n", i, c, c)
			switch c {
			case '\'', '"':
				if !inQuotes {
					inQuotes = true
					prevQuote = c
					if quoteIndex == 0 {
						p := strings.Fields(argString[:i])
						command.Args = append(command.Args, p...)
					}
					quoteIndex = i
				} else if c == prevQuote {
					inQuotes = false
					command.Args = append(command.Args, argString[quoteIndex+1:i])
					quoteIndex = i + 1
				}
			}
			prevC = c
		}
		_ = prevC
		if quoteIndex+1 < len(argString) {
			command.Args = append(command.Args, argString[quoteIndex+1:])
		}
	}

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command, nil
}

func (c *Command) Run() error {
	if c.BuiltinFunc != nil {
		if err := c.BuiltinFunc(c); err != nil {
			return err
		}

	} else {
		cmd := exec.Command(c.Path, c.Args...)
		cmd.Stdout = c.Stdout
		cmd.Stderr = c.Stderr
		// TODO: set the return value
		//if err := cmd.Run(); err != nil {
		_ = cmd.Run()
	}

	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stdout, "ERROR: failed to read: %s\n", err.Error())
			os.Exit(1)
		}

		command, err := ParseInput(input)
		if err != nil {
			fmt.Fprintf(os.Stdout, "%s\n", err)
			continue
		}

		if err := command.Run(); err != nil {
			fmt.Fprintf(os.Stdout, "%s\n", err)
		}
	}
}
