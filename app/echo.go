package main

import (
	"fmt"
	"os"
	"strings"
)

func HandleEcho(args []string) {
	var file string
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
		if item == ">" || item == "1>" {
			minx(&inF_i, i)
			inF = true
			if i+1 == len(args) {
				fmt.Println("echo: syntax error near unexpected token `newline'")
				return
			}
			file = args[i+1]
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
	if nl {
		oldStdout := os.Stdout
		if inF {
			f, err := os.Create(file)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			os.Stdout = f
		}
		fmt.Println(output)
		os.Stdout = oldStdout
	} else {
		oldStdout := os.Stdout
		if inF {
			f, err := os.Create(file)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			os.Stdout = f
		}
		fmt.Print(output)
		os.Stdout = oldStdout
	}
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
