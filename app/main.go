package main

import (
	"bufio"
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

func main() {
	// TODO: Uncomment the code below to pass the first stage
	comm := []string{"echo", "exit", "type"}
	for {
		fmt.Print("$ ")
		commandLn, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err == nil {
			if strings.ToLower(strings.TrimSpace(commandLn[:len(commandLn)-1])) == "exit" {
				break
			}
			command := strings.ToLower(strings.TrimSpace(strings.Split(commandLn[:len(commandLn)-1], " ")[0]))
			if command == "echo" {
				for i := 1; i < len(strings.Split(commandLn[:len(commandLn)-1], " ")); i++ {
					fmt.Print(strings.Split(commandLn[:len(commandLn)-1], " ")[i] + " ")
				}
				fmt.Println()
				continue
			}

			// Type command
			if command == "type" {
				comp := func(a, b string) bool {
					return a < b
				}
				// Built in
				arg := strings.ToLower(strings.TrimSpace(strings.Split(commandLn[:len(commandLn)-1], " ")[1]))
				if _, ch := BS(comm, arg, 0, 2, comp); ch {
					fmt.Printf("%s is a shell builtin", arg)
				} else {
					// Search for executable files using PATH.
					ch, path := Executable(arg)
					if ch {
						fmt.Printf("%s is %s", arg, path)
					} else {
						fmt.Printf("%s: not found", strings.ToLower(strings.TrimSpace(strings.Split(commandLn[:len(commandLn)-1], " ")[1])))
					}
				}
				fmt.Println()
				continue
			}

			// Not found
			fmt.Printf("%s: command not found", command)
			fmt.Println()
		} else {
			fmt.Print(err)
		}
	}
}
