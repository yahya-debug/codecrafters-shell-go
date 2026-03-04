package main

import (
	"fmt"
	"os"
	"os/exec"
)

func external_command(commandLn string) {
	arguments := ParseInput(commandLn)
	var outfile, errfile string
	var apnd bool
	tp := arguments[1:] // to pass arguments
	for i := 0; i < len(arguments); i++ {
		if arguments[i] == ">" || arguments[i] == "2>" || arguments[i] == ">>" || arguments[i] == "2>>" {
			arguments := ParseInput(commandLn)

			if i+1 >= len(arguments) {
				fmt.Println("syntax error near unexpected token `newline`")
				return
			}
			if arguments[i] == ">>" || arguments[i] == "2>>" {
				apnd = true
			}
			if arguments[i] == "2>" || arguments[i] == "2>>" {
				errfile = arguments[i+1]
			} else {
				outfile = arguments[i+1]
			}

			// remove ">" and filename from arguments
			tp = append(arguments[1:i], arguments[i+2:]...)
			break
		}
	}
	program := exec.Command(arguments[0], tp...)
	program.Stdin = os.Stdin
	program.Stdout = os.Stdout
	program.Stderr = os.Stderr
	// if there's a redirection print to the target file
	if outfile != "" {
		var f *os.File
		var err error
		if apnd {
			f, err = os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		} else {
			f, err = os.Create(outfile)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		program.Stdout = f
	}
	if errfile != "" {
		var f *os.File
		var err error
		if apnd {
			f, err = os.OpenFile(errfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		} else {
			f, err = os.Create(errfile)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		program.Stderr = f
	}
	// Print, read and report errors the terminal
	_ = program.Run()
}
