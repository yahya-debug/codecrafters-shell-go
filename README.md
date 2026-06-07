[![progress-banner](https://backend.codecrafters.io/progress/shell/d089bede-2e76-4cd1-92f4-92d29923363f)](https://app.codecrafters.io/users/yahya-debug?r=2qF)

# Build Your Own Shell — Go

A POSIX-flavored interactive shell built from scratch in Go as part of the [CodeCrafters "Build Your Own Shell"](https://app.codecrafters.io/courses/shell/overview) challenge. Every feature — from raw terminal input to pipeline execution — is hand-rolled with no third-party shell libraries.

---

## Demo

```
$ echo "hello world"
hello world
$ ls | grep go
go.mod
go.sum
$ cat file.txt > out.txt 2> err.txt
$ sleep 5 &
[1] 12345
$ jobs
[1]+  Running    sleep 5 &
$ history 3
  8  ls | grep go
  9  cat file.txt > out.txt 2> err.txt
  10  sleep 5 &
$ exit
```

---

## Features

### Core Shell
- [x] REPL loop with `$ ` prompt
- [x] Raw-mode terminal — full keystroke control via `golang.org/x/term`
- [x] Cursor-aware line editing (move, insert, delete at any position)
- [x] Bell feedback on invalid operations

### Input Parsing
- [x] Single-quoted strings (no interpretation)
- [x] Double-quoted strings with escape sequences
- [x] Backslash escaping outside quotes
- [x] Shell variable expansion: `$VAR` and `${VAR}`
- [x] `&&` sequential AND chaining
- [x] `&` background execution

### I/O
- [x] `|` multi-stage pipelines
- [x] `>` / `>>` stdout redirection (overwrite / append)
- [x] `2>` / `2>>` stderr redirection (overwrite / append)

### Built-in Commands
- [x] `echo [-n] [-e] [-E]` — with escape sequence interpretation
- [x] `cd [path|~]` — with home directory support
- [x] `pwd`
- [x] `type` — identifies builtins vs external executables with full path
- [x] `history [-r/-w/-a] [file]` — session history with file persistence
- [x] `jobs` — lists background jobs with `+`/`-` markers and Running/Done status
- [x] `complete -C <program> <command>` — programmable completion
- [x] `declare VAR=val / -p VAR` — shell variables
- [x] `exit` — flushes history to `$HISTFILE`

### Tab Completion
- [x] Complete commands (builtins + PATH executables)
- [x] Complete file/directory paths (appends `/` for dirs)
- [x] Longest-common-prefix completion on multiple matches
- [x] Double-Tab to list all candidates
- [x] Programmable completion via `-C` external programs (`COMP_LINE` / `COMP_POINT`)

### History
- [x] Up/Down arrow navigation
- [x] Load from `$HISTFILE` on startup, save on `exit`
- [x] `history -r/-w/-a` for manual file management

### Background Jobs
- [x] `cmd &` runs command asynchronously
- [x] Goroutine-based reaping via a `doneCh` channel
- [x] `jobs` with most-recent (`+`) and second-most-recent (`-`) markers

---

## Project Structure

```
app/
├── main.go          # REPL loop, built-in dispatch
├── ReadLine.go      # Raw terminal I/O, cursor control, tab completion
├── ParseInput.go    # Single-pass tokenizer: quotes, escapes, operators, variables
├── PipeLine.go      # Multi-stage pipe execution
├── ExternalComm.go  # External command runner + I/O redirection
├── jobs.go          # Background job control
├── History.go       # History navigation and file persistence
├── comlete.go       # Programmable completion (-C) support
├── cd.go            # cd builtin
├── echo.go          # echo builtin with -n/-e/-E flags
├── variables.go     # Shell variable store (declare)
└── yahya_library.go # Generic MergeSort, Binary Search, Pair[T,U]
```

---

## Running Locally

Requirements: **Go 1.25+**

```sh
./your_program.sh
```

Submit to CodeCrafters:

```sh
git push origin master
```
