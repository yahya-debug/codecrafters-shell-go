package main

import (
	"fmt"
	"os"
	"strings"
)

func HandleEcho(args []string) {
	var file, errfile string
	var nl, is, inF bool = true, false, false
	i := 0

	minx := func(a *int, b int) {
		if b < *a {
			*a = b
		}
	}
	inF_i := len(args)
	idc := 0
	for i < len(args) {
		item := args[i]
		valid := true
		if item == ">" || item == "1>" || item == "2>" {
			minx(&inF_i, i)
			inF = true
			if i+1 == len(args) {
				fmt.Println("echo: syntax error near unexpected token `newline'")
				return
			}
			if item == "2>" {
				errfile = args[i+1]
			} else {
				file = args[i+1]
			}
			i++
			continue
		} else if len(item) > 1 && item[0] == '-' && !inF && valid {
			for j := 1; j < len(item); j++ {
				switch item[j] {
				case 'n':
					nl = false
				case 'e':
					is = true
				case 'E':
					is = false
				default:
					valid, nl, is = false, true, false
				}
				if !valid {
					idc = 0
				} else {
					idc = i + 1
				}
			}
			i++
		} else {
			i++
		}
	}
	output := strings.Join(args[idc:inF_i], " ")
	if is {
		output = interpret(output)
	}
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	if file != "" {
		f, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		os.Stdout = f
	}

	if errfile != "" {
		f, err := os.Create(errfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		os.Stderr = f
	}
	if nl {
		fmt.Println(output)
	} else {
		fmt.Print(output)
	}

	os.Stdout = oldStdout
	os.Stderr = oldStderr
}

func interpret(str string) string {
	var res strings.Builder
	for c := 0; c < len(str); c++ {
		if str[c] == '\\' && c+1 < len(str) {
			c++
			switch str[c] {
			case 'n':
				res.WriteByte('\n')
			case 't':
				res.WriteByte('\t')
			case '\\':
				res.WriteByte('\\')
			case '"':
				res.WriteByte('"')
			case '\'':
				res.WriteByte('\'')
			default:
				res.WriteByte('\\')
				res.WriteByte(str[c+1])
			}
		} else {
			res.WriteByte(str[c])
		}
	}

	return res.String()
}
