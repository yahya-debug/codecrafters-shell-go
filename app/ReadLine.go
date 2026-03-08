package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const prompt = "$ "

var tabs int = 0
var cursor int = 0

var history []string
var histIndex int = -1
var l_append int = 0

func redraw(line []byte, toEnd bool) {
	fmt.Print("\r\033[K")
	fmt.Printf("%s%s", prompt, line)
	var cur int
	if cursor < len(line) {
		cur = len(line)
	} else {
		cursor = len(line)
	}
	if toEnd {
		for cursor < len(line) && cur < len(line) {
			fmt.Printf("%c%c%c", 27, 91, 67)
			cur++
		}
	} else {
		for cursor < cur && cur > 0 {
			fmt.Printf("%c%c%c", 27, 91, 68)
			cur--
		}
	}
}

func historyPrev() string {
	if len(history) == 0 || histIndex == -1 {
		return ""
	}
	if histIndex > 0 {
		histIndex--
	}
	return history[histIndex]
}
func historyNext() string {
	if len(history) == 0 {
		return ""
	}
	if histIndex < len(history)-1 {
		histIndex++
	} else {
		histIndex++
		return ""
	}
	return history[histIndex]
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
			cursor = 0
			if strings.TrimRight(string(line), "\r\n") != "" {
				history = append(history, strings.TrimRight(string(line), "\r\n"))
				histIndex = len(history)
			}
			return strings.TrimRight(string(line), "\r\n")
		case '\t':
			tabs++
			line = auto_complete(line)
			if tabs == 2 {
				tabs = 0
				redraw([]byte(line), true)
			}
			cursor = len(line)
			continue
		case 127:
			if cursor == 0 {
				fmt.Print("\a") // bell to ring
			} else if cursor > 0 {
				line = []byte(string(line[:cursor-1]) + string(line[cursor:]))
				cursor--
				redraw(line, false)
			}
		case 27: // ESC
			seq := make([]byte, 2)
			os.Stdin.Read(seq)

			if seq[0] == '[' {
				switch seq[1] {
				case 'A': // UP
					line = []byte(historyPrev())
					cursor = len(line)
					if strings.TrimSpace(string(line)) == "" {
						fmt.Print("\a") // bell to ring
					} else {
						redraw(line, true)
					}
				case 'B': // DOWN
					line = []byte(historyNext())
					cursor = len(line)
					if len(history) != 0 {
						redraw(line, true)
					}
				case 'C': // RIGHT
					if cursor == len(line) {
						fmt.Print("\a") // bell to ring
					} else {
						cursor++
						fmt.Printf("%c%c%c", buf[0], seq[0], seq[1])
					}

				case 'D': // LEFT
					if cursor == 0 {
						fmt.Print("\a") // bell to ring
					} else {
						cursor--
						fmt.Printf("%c%c%c", buf[0], seq[0], seq[1])
					}

				}
			}
		case 3: // Ctrl+C
			fmt.Println()
			fmt.Print("\r")
			return "exit"
		default:
			line = append(line[:cursor], append([]byte{buf[0]}, line[cursor:]...)...)
			cursor++
			fmt.Printf("%c", buf[0])
		}
		tabs = 0
	}
}

func find_fileMatching(str string) ([]Pair[string, bool], []string) {

	var dir string
	var prefix string

	lastSlash := strings.LastIndex(str, "/")

	if lastSlash == -1 || len(str) == 0 {
		dir = "."
		prefix = str
	} else {
		dir = str[:lastSlash+1]
		prefix = str[lastSlash+1:]
	}
	files, _ := os.ReadDir(dir)
	var matches []Pair[string, bool]
	var names []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) || str == "" {
			if lastSlash == -1 {
				pair := Pair[string, bool]{file.Name(), file.IsDir()}
				matches = append(matches, pair)
				names = append(names, file.Name())
			} else {
				pair := Pair[string, bool]{dir + file.Name(), file.IsDir()}
				matches = append(matches, pair)
				names = append(names, dir+file.Name())
			}
		}
	}
	return matches, names
}
func find_matching(str []byte) []string {
	var matches []string
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
	return matches
}

// TODO: Now we will find the longest common prefix among all matches (if exists)
func LCP(matches []string) string {
	if len(matches) == 0 {
		return ""
	}
	if len(matches) == 1 {
		return matches[0]
	}
	MergeSort(matches)
	it := 0
	for it < len(matches[0]) && it < len(matches[len(matches)-1]) && matches[0][it] == matches[len(matches)-1][it] {
		it++
	}
	return matches[0][:it]
}

var matches, f_names []string // map [number of matching prefixes] = {strings with this count}
var fileMatches []Pair[string, bool]

func auto_complete(str []byte) []byte {
	cmd := string(str)
	// complete commands
	completeComm := func() []byte {
		var ret strings.Builder // resulting string
		// first tab
		if tabs == 1 {

			matches = find_matching(str)
			lcp := LCP(matches)
			if len(matches) == 1 { // found 1 match in builtin commands
				l := matches[0] + " "
				cursor = len(l)
				redraw([]byte(l), true)
				for i := 0; i < len(matches[0]); i++ {
					ret.WriteByte(matches[0][i])
				}
				ret.WriteByte(' ')
			} else if len(matches) > 1 && len(lcp) > len(cmd) {
				for i := 0; i < len(lcp); i++ {
					ret.WriteByte(lcp[i])
				}
				redraw([]byte(ret.String()), true)
			} else {
				for i := 0; i < len(cmd); i++ {
					ret.WriteByte(cmd[i])
				}
				fmt.Print("\a") // bell to ring
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
			} else {
				for i := 0; i < len(cmd); i++ {
					ret.WriteByte(cmd[i])
				}
				fmt.Print("\a") // bell to ring
			}
		}
		return []byte(ret.String())
	}

	// complete for file names
	completeFile := func() []byte {
		var ret strings.Builder // resulting string
		cmd_parts := ParseInput(strings.TrimLeft(cmd, " "))
		if tabs == 1 {
			cur_F := cmd_parts[len(cmd_parts)-1]
			if cmd[len(cmd)-1] == ' ' {
				cur_F = ""
			}
			cut := strings.TrimSuffix(cmd, cur_F)
			for i := 0; i < len(cut); i++ {
				ret.WriteByte(cut[i])
			}
			fileMatches, f_names = find_fileMatching(cur_F)
			lcp := LCP(f_names)
			if len(fileMatches) == 1 {
				var l string
				for i := 0; i < len(fileMatches[0].f); i++ {
					ret.WriteByte(fileMatches[0].f[i])
				}
				l = strings.TrimSuffix(cmd, cur_F)
				if fileMatches[0].s {
					l += fileMatches[0].f + "/"
					ret.WriteByte('/')
				} else {
					l += fileMatches[0].f + " "
					ret.WriteByte(' ')
				}
				cursor = len(l)
				redraw([]byte(l), true)
				tabs = 0
			} else if len(fileMatches) > 1 && len(lcp) > len(cur_F) {
				for i := 0; i < len(lcp); i++ {
					ret.WriteByte(lcp[i])
				}
				redraw([]byte(ret.String()), true)
			} else {
				ret.Reset()
				for i := 0; i < len(cmd); i++ {
					ret.WriteByte(cmd[i])
				}
				fmt.Print("\a") // bell to ring
			}
		} else {

			for i := 0; i < len(cmd); i++ {
				ret.WriteByte(cmd[i])
			}
			if len(fileMatches) > 1 {
				fmt.Println("\r")
				for _, i := range fileMatches {
					if i.s {
						fmt.Printf("%s/  ", i.f)
					} else {
						fmt.Printf("%s  ", i.f)

					}
				}
				fmt.Println("\r")
			} else {
				fmt.Print("\a") // bell to ring
			}
		}

		return []byte(ret.String())
	}
	if !strings.Contains(strings.TrimLeft(cmd, " "), " ") {
		return completeComm()
	}
	return completeFile()
}
