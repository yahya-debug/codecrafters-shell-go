package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func ReadLine() string {
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	buf := make([]byte, 1)
	var line []byte
	for {
		os.Stdin.Read(buf)
		switch buf[0] {
		case '\n', '\r':
			fmt.Println()
			fmt.Print("\r")
			return strings.TrimRight(string(line), "\r\n")
		case '\t':
			line = auto_complete(line)
			continue
		case 127:
			if len(line) > 0 {
				line = line[:len(line)-1]
				fmt.Print("\b \b")
			}
		case 3: // Ctrl+C
			fmt.Println()
			fmt.Print("\r")
			return "exit"
		default:
			line = append(line, buf[0])
			fmt.Printf("%c", buf[0])
		}
	}
}

func auto_complete(str []byte) []byte {
	cmd := string(str)
	var ret strings.Builder
	matches := make(map[int][]string)
	mx := 0
	for _, com := range comm {
		for i := 0; i < len(cmd); i++ {
			if cmd[i] != com[i] {
				maxx(&mx, i)
				matches[i] = append(matches[i], com)
				break
			}
			if i == len(cmd)-1 {
				maxx(&mx, i+1)
				matches[i+1] = append(matches[i+1], com)
			}
		}
	}
	if mx > 0 && len(matches[mx]) == 1 {
		fmt.Printf("\r$ %s ", matches[mx][0])
		for i := 0; i < len(matches[mx][0]); i++ {
			ret.WriteByte(matches[mx][0][i])
		}
		ret.WriteByte(' ')
	} else {
		for i := 0; i < len(cmd); i++ {
			ret.WriteByte(cmd[i])
		}
	}
	return []byte(ret.String())
}
