package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage
	for {
		fmt.Print("$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err == nil {
			if strings.EqualFold(strings.ToLower(strings.TrimSpace(command[:len(command)-1])), "exit") {
				break
			} else if strings.ToLower(strings.TrimSpace(strings.Split(command[:len(command)-1], " ")[0])) == "echo" {
				for i := 1; i < len(strings.Split(command[:len(command)-1], " ")); i++ {
					fmt.Println(strings.Split(command[:len(command)-1], " ")[i] + " ")
				}
				continue
			}
			fmt.Printf("%s: command not found", command[:len(command)-1])
			fmt.Println()
		} else {
			fmt.Print(err)
		}
	}
}
