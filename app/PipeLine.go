package main

import (
	"os"
	"os/exec"
	"strings"
)

func runPipeline(commands ...[]string) {

	n := len(commands)

	var prevReader *os.File
	var procs []*exec.Cmd

	for i := 0; i < n; i++ {

		var in *os.File
		var out *os.File

		if i == 0 {
			in = os.Stdin
		} else {
			in = prevReader
		}

		if i == n-1 {
			out = os.Stdout
		} else {
			r, w, _ := os.Pipe()
			out = w
			prevReader = r
		}

		cmdName := strings.TrimSpace(commands[i][0])

		// ---------- BUILTINS ----------
		if cmdName == "echo" {

			old := os.Stdout
			os.Stdout = out
			HandleEcho(commands[i][1:])
			os.Stdout = old

			if out != os.Stdout {
				out.Close()
			}

			continue
		}

		if cmdName == "pwd" {

			old := os.Stdout
			os.Stdout = out

			if pwd, err := os.Getwd(); err == nil {
				println(pwd)
			}

			os.Stdout = old

			if out != os.Stdout {
				out.Close()
			}

			continue
		}
		if cmdName == "type" {

			old := os.Stdout
			os.Stdout = out

			comp := func(a, b string) bool {
				return a < b
			}

			for j := 1; j < len(commands[i]); j++ {

				arg := strings.TrimSpace(commands[i][j])

				if _, ch := BS(comm, arg, 0, len(comm)-1, comp); ch {
					println(arg + " is a shell builtin")
				} else {
					ch, path := Executable(arg)
					if ch {
						println(arg + " is " + path)
					} else {
						println(arg + ": not found")
					}
				}
			}

			os.Stdout = old

			if out != os.Stdout {
				out.Close()
			}

			continue
		}
		// ---------- EXTERNAL ----------
		cmd := exec.Command(commands[i][0], commands[i][1:]...)

		cmd.Stdin = in
		cmd.Stdout = out
		cmd.Stderr = os.Stderr

		cmd.Start()

		procs = append(procs, cmd)

		if out != os.Stdout {
			out.Close()
		}
	}

	// wait for external commands
	for _, p := range procs {
		p.Wait()
	}
}
