package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

// Binary search
func BS[T any](arr []T, target T, l int, r int, less func(a, b T) bool) (int, bool) {
	if r < l {
		return -1, false
	}
	mid := (l + r) / 2
	if less(arr[mid], target) {
		return BS(arr, target, mid+1, r, less)
	} else if less(target, arr[mid]) {
		return BS(arr, target, l, mid-1, less)
	}
	return mid, true
}

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

var comm []string = []string{"echo", "exit", "pwd", "type"}

func main() {
	// TODO: Uncomment the code below to pass the first stage
	for {
		fmt.Print("$ ")
		commandLn := ReadLine()
		if commandLn == "" {
			continue
		}
		command := ParseInput(commandLn)[0]
		if command == "exit" || commandLn == "" {
			break
		}
		if command == "echo" {
			if idx := strings.Index(commandLn, " "); idx != -1 {
				HandleEcho(ParseInput(commandLn[idx+1:]))
			}
			continue
		}

		// Type command
		if command == "type" {
			comp := func(a, b string) bool {
				return a < b
			}
			// Built in
			arg := strings.ToLower(strings.TrimSpace(strings.Split(commandLn, " ")[1]))
			if _, ch := BS(comm, arg, 0, len(comm)-1, comp); ch {
				fmt.Printf("%s is a shell builtin", arg)
			} else {
				// Search for executable files using PATH.
				ch, path := Executable(arg)
				if ch {
					fmt.Printf("%s is %s", arg, path)
				} else {
					fmt.Printf("%s: not found", strings.ToLower(strings.TrimSpace(strings.Split(commandLn, " ")[1])))
				}
			}
			fmt.Println()
			continue
		}
		// get working directory
		if command == "pwd" {
			if pwd, err := os.Getwd(); err == nil {
				fmt.Println(pwd)
			}
			continue
		}
		// Handle absolute path
		if command == "cd" {
			arg := strings.TrimSpace(strings.Split(commandLn, " ")[1])
			d := HandleCD(arg)
			if !d {
				fmt.Printf("\rcd: %s: No such file or directory\n", arg)
			}
			continue
		}
		// Run external command
		if ok, _ := Executable(command); ok {
			external_command(commandLn)
			continue
		}
		// Not found
		fmt.Printf("\r%s: command not found\n", command)
	}
}
