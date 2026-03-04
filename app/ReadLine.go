package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const prompt = "$ "

var tabs int = 0

func redraw(line []byte) {
	fmt.Print("\r\033[K")
	fmt.Printf("%s%s", prompt, line)
}

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
			tabs = 0
			return strings.TrimRight(string(line), "\r\n")
		case '\t':
			tabs++
			line = auto_complete(line)
			if tabs == 2 {
				tabs = 0
				return ""
			}
			continue
		case 127:
			if len(line) > 0 {
				line = line[:len(line)-1]
				redraw(line)
			}
		case 3: // Ctrl+C
			fmt.Println()
			fmt.Print("\r")
			return "exit"
		default:
			line = append(line, buf[0])
			fmt.Printf("%c", buf[0])
		}
		tabs = 0
	}
}

func auto_complete(str []byte) []byte {
	// cast to string
	cmd := string(str)
	var ret strings.Builder // resulting string
	var matches []string    // map [number of matching prefixes] = {strings with this count}
	for _, com := range comm {
		if strings.HasPrefix(com, cmd) {
			matches = append(matches, com)
		}
	}
	if len(matches) == 1 {
		l := matches[0] + " "
		redraw([]byte(l))
		for i := 0; i < len(matches[0]); i++ {
			ret.WriteByte(matches[0][i])
		}
		ret.WriteByte(' ')
	} else if len(matches) > 1 && tabs == 2 {
		fmt.Println("\r")
		for _, i := range matches {
			fmt.Printf("%s	", i)
		}
		fmt.Println("\r")
	} else {
		d := BSs(execs, cmd, 0, len(execs)-1, func(a, b string) bool {
			return a < b
		})

		for i := d; i < len(execs); i++ {
			com := execs[i]
			if strings.HasPrefix(com, cmd) {
				matches = append(matches, com)
			}
		}
		if len(matches) == 1 {
			l := matches[0] + " "
			redraw([]byte(l))
			for i := 0; i < len(matches[0]); i++ {
				ret.WriteByte(matches[0][i])
			}
			ret.WriteByte(' ')
		} else if len(matches) > 1 && tabs == 2 {
			fmt.Println("\r")
			for _, i := range matches {
				fmt.Printf("%s	", i)
			}
			fmt.Println("\r")
		} else {
			fmt.Print("\a")
			for i := 0; i < len(cmd); i++ {
				ret.WriteByte(cmd[i])
			}
		}
	}
	return []byte(ret.String())
}
