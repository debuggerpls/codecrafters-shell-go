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

## The cd builtin: absolute paths
```
go run cmd/myshell/main.go 
$ cd /usr/local/bin
$ pwd
/usr/local/bin
$ cd /does_not_exist
cd: /does_not_exist: No such file or directory
$
```

## The cd builtin: relative paths
```
go run cmd/myshell/main.go 
$ cd /usr
$ pwd
/usr
$ cd ./local/bin
$ pwd
/usr/local/bin
$ cd ../../
$ pwd
/usr
$
```

## The cd builtin: home directory
```
go run cmd/myshell/main.go 
$ cd /usr/local/bin
$ pwd
/usr/local/bin
$ cd ~
$ pwd
/home/user
$
```

## Single quotes
```
go run cmd/myshell/main.go 
$ echo 'shell hello'
shell hello
$ echo 'world     test'
world     test
$
```

## Double quotes
```
go run cmd/myshell/main.go 
$ echo "quz  hello"  "bar"
quz  hello bar
$ echo "bar"  "shell's"  "foo"
bar shell's foo
$
```

# TODO FROM HERE
## Backslash outside quotes
https://www.gnu.org/software/bash/manual/bash.html#Escape-Character
```
go run cmd/myshell/main.go 
$ echo "before\   after"
before\   after
$ echo world\ \ \ \ \ \ script
world      script
$
$ cat "/tmp/file\\name" "/tmp/file\ name" 
content1 content2
```

## Backslash within single quotes
https://www.gnu.org/software/bash/manual/bash.html#Single-Quotes
```
go run cmd/myshell/main.go 
$ echo 'shell\\\nscript'
shell\\\nscript
$ echo 'example\"testhello\"shell'
example\"testhello\"shell
$
$ cat "/tmp/file/'name'" "/tmp/file/'\name\'"  
content1 content2
```

## Backslash within double quotes
https://www.gnu.org/software/bash/manual/bash.html#Double-Quotes
```
go run cmd/myshell/main.go 
$ echo "hello'script'\\n'world"
hello'script'\n'world
$ echo "hello\"insidequotes"script\"
hello"insidequotesscript"
$
$ cat "/tmp/"file\name"" "/tmp/"file name"" 
content1 content2
```

## Executing a quoted executable
```
go run cmd/myshell/main.go 
$ 'exe with "quotes"' file
content1
$ "exe with 'single quotes'" file
content2
```

## Redirect stdout
https://www.gnu.org/software/bash/manual/bash.html#Redirecting-Output
```
go run cmd/myshell/main.go 
$ ls /tmp/baz > /tmp/foo/baz.md
$ cat /tmp/foo/baz.md
apple
blueberry
$ echo 'Hello James' 1> /tmp/foo/foo.md
$ cat /tmp/foo/foo.md
Hello James
$ cat /tmp/baz/blueberry nonexistent 1> /tmp/foo/quz.md
cat: nonexistent: No such file or directory
$ cat /tmp/foo/quz.md
blueberry
```

## Redirect stderr
https://www.gnu.org/software/bash/manual/bash.html#Redirecting-Output
```
go run cmd/myshell/main.go 
$ ls nonexistent 2> /tmp/quz/baz.md
$ cat /tmp/quz/baz.md
ls: cannot access 'nonexistent': No such file or directory
$ echo 'Maria file cannot be found' 2> /tmp/quz/foo.md
Maria file cannot be found
$ cat /tmp/bar/pear nonexistent 2> /tmp/quz/quz.md
pear
$ cat /tmp/quz/quz.md
cat: nonexistent: No such file or directory
```

## Append stdout
https://www.gnu.org/software/bash/manual/bash.html#Appending-Redirected-Output
```
go run cmd/myshell/main.go 
$ ls /tmp/baz >> /tmp/bar/bar.md
$ cat /tmp/bar/bar.md
apple
banana
blueberry
$ echo 'Hello Emily' 1>> /tmp/bar/baz.md
$ echo 'Hello Maria' 1>> /tmp/bar/baz.md
$ cat /tmp/bar/baz.md
Hello Emily
Hello Maria
$ echo "List of files: " > /tmp/bar/qux.md
$ ls /tmp/baz >> /tmp/bar/qux.md
$ cat /tmp/bar/qux.md
List of files:
apple
banana
blueberry
```

## Append stderr
https://www.gnu.org/software/bash/manual/bash.html#Appending-Redirected-Output
```
go run cmd/myshell/main.go 
$ ls nonexistent >> /tmp/foo/baz.md
ls: cannot access 'nonexistent': No such file or directory
$ ls nonexistent 2>> /tmp/foo/qux.md
$ cat /tmp/foo/qux.md
ls: cannot access 'nonexistent': No such file or directory
$ echo "James says Error" 2>> /tmp/foo/quz.md
James says Error
$ cat nonexistent 2>> /tmp/foo/quz.md
$ ls nonexistent 2>> /tmp/foo/quz.md
$ cat /tmp/foo/quz.md
cat: nonexistent: No such file or directory
ls: cannot access 'nonexistent': No such file or directory
```
