[![progress-banner](https://backend.codecrafters.io/progress/shell/4263e779-28ac-4378-b2a1-be658c028680)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This is a starting point for Go solutions to the
["Build Your Own Shell" Challenge](https://app.codecrafters.io/courses/shell/overview).

In this challenge, you'll build your own POSIX compliant shell that's capable of
interpreting shell commands, running external programs and builtin commands like
cd, pwd, echo and more. Along the way, you'll learn about shell command parsing,
REPLs, builtin commands, and more.

## Invalid commands
```
go run cmd/myshell/main.go
$ hello
hello: command not found
```

## REPL
```
go run cmd/myshell/main.go
$ hello
hello: command not found
$ hello
hello: command not found
$
```

## The exit builtin
```
go run cmd/myshell/main.go 
$ hello
hello: command not found
$ exit
```

## The echo builtin
```
go run cmd/myshell/main.go 
$ echo hello world
hello world
$ exit
```

## The type builtin
```
go run cmd/myshell/main.go 
$ type invalid
invalid: not found
$ type echo
echo is a shell builtin
$ exit
```

## The type builtin: executable files
```
PATH="/usr/bin:/usr/local/bin" ./your_program.sh
$ type ls
ls is /usr/bin/ls
$ type abcd
abcd is /usr/local/bin/abcd
$ type invalid_command
invalid_command: not found
$
```

## Run a program
```
go run cmd/myshell/main.go 
$ ls
cmd  codecrafters.yml  go.mod  go.sum  README.md  your_program.sh
$ 
```

## The pwd builtin
```
go run cmd/myshell/main.go 
$ pwd
/home/debuggerpls/code/codecrafters-shell-go
$ type pwd
pwd is a shell builtin
```
