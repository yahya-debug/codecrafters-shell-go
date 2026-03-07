package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getExecs() {
	pathEnv := os.Getenv("PATH")
	dirs := strings.Split(pathEnv, ":")
	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, f := range files {
			info, err := f.Info()
			if err != nil {
				continue
			}
			if !info.IsDir() && info.Mode()&0111 != 0 {
				execs = append(execs, f.Name())
			}
		}
	}
}

func external_command(commandLn []string, in, out, errOut *os.File) {

	var outfile, errfile string
	var apnd bool
	tp := commandLn[1:]

	for i := 0; i < len(commandLn); i++ {

		if commandLn[i] == ">" || commandLn[i] == "2>" || commandLn[i] == ">>" || commandLn[i] == "2>>" {

			if i+1 >= len(commandLn) {
				fmt.Println("syntax error near unexpected token `newline`")
				return
			}

			if commandLn[i] == ">>" || commandLn[i] == "2>>" {
				apnd = true
			}

			if commandLn[i] == "2>" || commandLn[i] == "2>>" {
				errfile = commandLn[i+1]
			} else {
				outfile = commandLn[i+1]
			}

			tp = append(commandLn[1:i], commandLn[i+2:]...)
			break
		}
	}

	program := exec.Command(commandLn[0], tp...)

	// use pipe/terminal passed from caller
	program.Stdin = in
	program.Stdout = out
	program.Stderr = errOut

	// handle stdout redirection
	if outfile != "" && out == os.Stdout {

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

	// handle stderr redirection
	if errfile != "" && errOut == os.Stderr {

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

	_ = program.Run()
}
