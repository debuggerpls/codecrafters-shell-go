package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FIXME: there must be a better solution to overcome initialization loop in typeBuiltin
var builtinsNames = []string{
	"exit",
	"echo",
	"type",
	"pwd",
	"cd",
}

var builtins = map[string]func([]string){
	builtinsNames[0]: exitBuiltin,
	builtinsNames[1]: echoBuiltin,
	builtinsNames[2]: typeBuiltin,
	builtinsNames[3]: pwdBuiltin,
	builtinsNames[4]: cdBuiltin,
}

func IsBuiltin(name string) bool {
	for _, builtin := range builtinsNames {
		if builtin == name {
			return true
		}
	}
	return false
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

func exitBuiltin(args []string) {
	// TODO: is it ok to just exit here?
	// TODO: pass exit arg when exiting
	os.Exit(0)
}

func echoBuiltin(args []string) {
	if len(args) == 0 {
		return
	}

	fmt.Fprintf(os.Stdout, "%s\n", strings.Join(args, " "))
}

func typeBuiltin(args []string) {
	if len(args) == 0 {
		return
	}

	command := args[0]
	if IsBuiltin(command) {
		fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", command)
		return
	}
	if path, ok := GetPath(command); ok {
		fmt.Fprintf(os.Stdout, "%s is %s\n", command, path)
		return
	}

	fmt.Fprintf(os.Stdout, "%s: not found\n", command)
}

func pwdBuiltin(args []string) {
	//fmt.Fprintf(os.Stdout, "%s\n", os.Getenv("PWD"))
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	fmt.Fprintf(os.Stdout, "%s\n", pwd)
}

func cdBuiltin(args []string) {
	var dir string
	if len(args) == 0 || args[0] == "~" {
		dir = os.Getenv("HOME")
	} else {
		dir = args[0]
	}

	if err := os.Chdir(dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", args[0])
			return
		} else {
			// TODO: do not pass the underlying error
			fmt.Fprintf(os.Stderr, "cd: %s: %s\n", args[0], err)
		}
	}
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

		inputParts := strings.Split(strings.TrimSpace(input), " ")
		command := inputParts[0]

		// process builtins
		builtin, ok := builtins[command]
		if ok {
			builtin(inputParts[1:])
			continue
		}

		// process executables
		fpath, ok := GetPath(command)
		if ok {
			cmd := exec.Command(fpath, inputParts[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			// TODO: set the return value
			//if err := cmd.Run(); err != nil {
			_ = cmd.Run()
			continue
		}

		fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
	}
}
