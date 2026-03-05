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
				redraw([]byte(line))
				// line = append(line, buf[0])
			}
			// fmt.Printf("\n %s \n", string(line))
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

var matches []string // map [number of matching prefixes] = {strings with this count}
func find_matching(str []byte) {
	// cast to string
	cmd := string(str)

	// Find matches in builtin commands
	for _, com := range comm {
		if strings.HasPrefix(com, cmd) {
			matches = append(matches, com)
		}
	}

	// if no matches in the built in commands
	if len(matches) == 0 {
		// search in executables
		// use binary search to fastly find the matching sequence
		d := BSs(execs, cmd, 0, len(execs)-1, func(a, b string) bool {
			return a < b
		})

		for i := d; i < len(execs); i++ {
			com := execs[i]
			if strings.HasPrefix(com, cmd) {
				matches = append(matches, com)
			}
		}
	}
}

// TODO: Now we will find the longest common prefix among all matches (if exists)
func LCP() string {
	if len(matches) == 0 {
		return ""
	}
	if len(matches) == 1 {
		return matches[0]
	}
	MergeSort(matches)
	it := 0
	for it < len(matches[0]) && it < len(matches[len(matches)-1]) && matches[0][i] == matches[len(matches)-1][i] {
		it++
	}
	return matches[0][:it]
}
func auto_complete(str []byte) []byte {
	cmd := string(str)

	var ret strings.Builder // resulting string
	// first tab
	if tabs == 1 {
		clear(matches)
		find_matching(str)
		if len(matches) == 1 { // found 1 match in builtin commands
			l := matches[0] + " "
			redraw([]byte(l))
			for i := 0; i < len(matches[0]); i++ {
				ret.WriteByte(matches[0][i])
			}
			ret.WriteByte(' ')
		} else if len(matches) > 1 {
			lcp := LCP()
			fmt.Println("\r")
			for i := 0; i < len(lcp); i++ {
				ret.WriteByte(lcp[i])
			}
		}
	} else {
		if len(matches) > 1 {
			fmt.Println("\r")
			for _, i := range matches {
				fmt.Printf("%s	", i)
			}
			fmt.Println("\r")
			for i := 0; i < len(cmd); i++ {
				ret.WriteByte(cmd[i])
			}
		}
	}

	// if len(matches) > 1 && tabs == 2 { // many in builtin commands -> handle duplicated tabs
	// 	fmt.Println("\r")
	// 	for _, i := range matches {
	// 		fmt.Printf("%s	", i)
	// 	}
	// 	fmt.Println("\r")
	// 	for i := 0; i < len(cmd); i++ {
	// 		ret.WriteByte(cmd[i])
	// 	}
	// } else {
	// 	// if multiple but 1 tab || no match -> bell to ring
	// 	fmt.Print("\a") // bell to ring
	// 	for i := 0; i < len(cmd); i++ {
	// 		ret.WriteByte(cmd[i])
	// 	}
	// }
	return []byte(ret.String())
}
