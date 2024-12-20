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
