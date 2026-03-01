package main

import (
	"bufio"
	"fmt"
	"os"
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

func main() {
	// TODO: Uncomment the code below to pass the first stage
	comm := []string{"echo", "exit", "type"}
	for {
		fmt.Print("$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err == nil {
			if strings.EqualFold(strings.ToLower(strings.TrimSpace(command[:len(command)-1])), "exit") {
				break
			}
			if strings.ToLower(strings.TrimSpace(strings.Split(command[:len(command)-1], " ")[0])) == "echo" {
				for i := 1; i < len(strings.Split(command[:len(command)-1], " ")); i++ {
					fmt.Print(strings.Split(command[:len(command)-1], " ")[i] + " ")
				}
				fmt.Println()
				continue
			}
			if strings.ToLower(strings.TrimSpace(strings.Split(command[:len(command)-1], " ")[0])) == "type" {
				comp := func(a, b string) bool {
					return a < b
				}
				if _, ch := BS(comm, strings.ToLower(strings.TrimSpace(strings.Split(command[:len(command)-1], " ")[1])), 0, 2, comp); ch {
					fmt.Printf("%s is a shell builtin", strings.ToLower(strings.TrimSpace(strings.Split(command[:len(command)-1], " ")[1])))
				} else {
					fmt.Printf("%s: not found", strings.ToLower(strings.TrimSpace(strings.Split(command[:len(command)-1], " ")[1])))
				}
				fmt.Println()
				continue
			}
			fmt.Printf("%s: command not found", command[:len(command)-1])
			fmt.Println()
		} else {
			fmt.Print(err)
		}
	}
}
