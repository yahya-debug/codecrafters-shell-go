package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func Executable(arg string) (bool, string) {
	pathEnv := os.Getenv("PATH")
	dirs := strings.Split(pathEnv, ":")
	for _, dir := range dirs {
		full_path := filepath.Join(dir, arg)
		info, err := os.Stat(full_path)
		if err != nil {
			continue
		}
		if !info.IsDir() && info.Mode()&0111 != 0 {
			return true, full_path
		}
	}
	return false, ""
}

var comm []string = []string{"cd", "echo", "exit", "pwd", "type"}
var execs []string

func main() {
	// Get Executable eternal commands and sort them
	// sorting will reduce time and will also help us print them in the way normal shell does
	getExecs()
	execs = MergeSort(execs)

	// run input and wait for user to press enter
	for {
		fmt.Print("$ ")
		commandLn := ReadLine()
		if commandLn == "" {
			continue
		}
		inp := ParseInput(commandLn)
		// if
		command := inp[0]
		if command == "exit" || commandLn == "" {
			break
		}
		var args [][]string
		l := 0
		for i := 0; i < len(inp); i++ {
			if inp[i] == "|" && i > 0 {
				args = append(args, inp[l:i])
				l = i + 1
			}
		}
		if inp[len(inp)-1] != "|" {
			args = append(args, inp[l:])
		}
		fmt.Print(run(args...))
	}
}

// here we run commands
func run(commands ...[]string) string {
	var out string
	if len(commands) > 1 {
		runPipeline(commands...)
		return ""
	}

	for i := 0; i < len(commands); i++ {
		command := commands[i][0]
		// Type command
		if command == "echo" {
			HandleEcho(commands[i][1:])
			continue
		}
		if command == "type" {
			comp := func(a, b string) bool {
				return a < b
			}
			// Built in
			for j := 1; j < len(commands[i]); j++ {
				arg := strings.TrimSpace(commands[i][j])
				if _, ch := BS(comm, arg, 0, len(comm)-1, comp); ch {
					out += arg + " is a shell builtin"
				} else {
					// Search for executable files using PATH.
					ch, path := Executable(arg)
					if ch {
						out += arg + " is " + path
					} else {
						out += commands[i][j] + ": not found"
					}
				}
				out += "\n"
			}
			continue
		}
		// get working directory
		if command == "pwd" {
			if pwd, err := os.Getwd(); err == nil {
				out += pwd + "\n"
			}
			continue
		}
		// Handle absolute path
		if command == "cd" {
			arg := ParseInput(strings.Join(commands[i], " "))[1]
			d := HandleCD(arg)
			if !d {
				out += "cd: + " + arg + ": No such file or directory\n"
			}
			continue
		}
		// Run external command
		if ok, _ := Executable(command); ok {
			external_command(strings.Join(commands[i], " "), os.Stdin, os.Stdout, os.Stderr)
			continue
		}
		// Not found
		out += command + ": command not found\n"
	}
	return out
}
