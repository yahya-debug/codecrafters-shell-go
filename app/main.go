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

var comm []string = []string{"echo", "exit", "pwd", "type"}
var execs []string

func main() {
	getExecs := func() {
		pathEnv := os.Getenv("PATH")
		dirs := strings.Split(pathEnv, ":")
		for _, dir := range dirs {
			files, err := os.ReadDir(dir)
			if err != nil {
				continue
			}
			for _, f := range files {
				info, err := f.Info()
				if err != nil {
					continue
				}
				if !info.IsDir() && info.Mode()&0111 != 0 {
					execs = append(execs, f.Name())
				}
			}
		}
	}
	getExecs()
	execs = MergeSort(execs)
	// for _, i := range execs {
	// 	fmt.Printf("%s ", i)
	// }
	// fmt.Println()
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
