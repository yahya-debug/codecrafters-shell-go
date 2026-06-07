package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func Executable(arg string) (bool, string) {
	pathEnv := os.Getenv("PATH")
	dirs := strings.Split(pathEnv, ":")
	for _, dir := range dirs {
		full_path := filepath.Join(dir, arg)
		info, err := os.Stat(full_path)
		if err != nil {
			continue
		}
		if !info.IsDir() && info.Mode()&0111 != 0 {
			return true, full_path
		}
	}
	return false, ""
}

var comm []string = []string{"cd", "complete", "declare", "echo", "exit", "history", "jobs", "pwd", "type"}
var execs []string
var hist_def_file = os.Getenv("HISTFILE")

func main() {
	// Set Defaults
	if hist_def_file != "" {
		ReadHist(hist_def_file)
	}

	// Get Executable eternal commands and sort them
	// sorting will reduce time and will also help us print them in the way normal shell does
	getExecs()
	execs = MergeSort(execs)
	alive := true
	// run input and wait for user to press enter
	for {
		if !alive {
			break
		}
		fmt.Print("$ ")
		commandLn := ReadLine()
		if commandLn == "" {
			continue
		}
		inp, multi := ParseInput(commandLn)
		if multi {
			if inp[len(inp)-1] != "&" {
				var inps [][]string
				l := 0
				for i := 0; i < len(inp); i++ {
					if inp[i] == "&&" {
						inps = append(inps, inp[l:i])
						l = i + 1
					}
				}
				inps = append(inps, inp[l:])
				for i := 0; i < len(inps); i++ {
					var success bool
					alive, success = processLine(inps[i])
					if !alive || !success {
						break
					}
				}
				continue
			}
		}
		alive, _ = processLine(inp)
	}
}

func processLine(inp []string) (bool, bool) {
	command := inp[0]
	if command == "exit" {
		WriteHist(hist_def_file)
		return false, true
	}

	if inp[len(inp)-1] == "&" {
		jobRun(inp[:len(inp)-1])
		return true, true
	}
	var args [][]string
	l := 0
	for i := 0; i < len(inp); i++ {
		if inp[i] == "|" && i > 0 {
			args = append(args, inp[l:i])
			l = i + 1
		}
	}
	if inp[len(inp)-1] != "|" {
		args = append(args, inp[l:])
	}
	out, success := run(args...)
	fmt.Print("\r" + out)
	return true, success
}

// here we run commands
func run(commands ...[]string) (string, bool) {
	var out string
	if len(commands) > 1 {
		runPipeline(commands...)
		return "", true
	}
	success := true

	for i := 0; i < len(commands); i++ {
		command := strings.TrimSpace(commands[i][0])
		// Type command
		if command == "echo" {
			HandleEcho(commands[i][1:])
			continue
		}
		if command == "type" {
			comp := func(a, b string) bool {
				return a < b
			}
			// Built in
			for j := 1; j < len(commands[i]); j++ {
				arg := strings.TrimSpace(commands[i][j])
				if _, ch := BS(comm, arg, 0, len(comm)-1, comp); ch {
					out += arg + " is a shell builtin"
				} else {
					// Search for executable files using PATH.
					ch, path := Executable(arg)
					if ch {
						out += arg + " is " + path
					} else {
						out += commands[i][j] + ": not found"
					}
				}
				out += "\n"
			}
			continue
		}

		// History
		if command == "history" {
			var err error
			if len(commands[i]) > 1 {
				_, err = strconv.Atoi(commands[i][1])
			} else {
				err = nil
			}
			if len(commands[i]) <= 2 && err == nil {
				var i int
				if len(commands[i]) == 1 { // Deafault -> print all history
					i = len(history)
				} else { // If user specified a number to print history items as much as it
					i, _ = strconv.Atoi(commands[i][1])
				}
				for i = len(history) - i; i < len(history); i++ {
					if i >= 0 {
						out += "  " + strconv.Itoa(i+1) + "  " + history[i] + "\n"
					}
				}
			} else {
				ch := commands[i][1]
				switch ch {
				case "-r":
					if len(commands[i]) >= 3 {
						ReadHist(commands[i][2])
					}
				case "-w":
					if len(commands[i]) < 3 {
						continue
					}
					WriteHist(commands[i][2])
				case "-a":
					if len(commands[i]) < 3 {
						continue
					}
					file, err := os.OpenFile(commands[i][2], os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
					if err != nil {
						continue
					}
					defer file.Close()
					file_writer := bufio.NewWriter(file)
					for i := l_append; i < len(history); i++ {
						file_writer.WriteString(history[i] + "\n")
					}
					l_append = len(history)
					file_writer.Flush()
				}
			}
			continue
		}

		// get working directory
		if command == "pwd" {
			if pwd, err := os.Getwd(); err == nil {
				out += pwd + "\n"
			}
			continue
		}

		// Handle absolute path
		if command == "cd" {
			arg_, _ := ParseInput(strings.Join(commands[i], " "))
			arg := arg_[1]
			if !HandleCD(arg) {
				out += "cd: " + arg + ": No such file or directory\n"
				success = false
			}
			continue
		}

		// jobs
		if command == "jobs" {
			allJobs := showJobs(false)
			for i := len(allJobs) - 1; i >= 0; i-- {
				out += allJobs[i]
			}
			continue
		}

		// complete
		if command == "complete" {
			// complete -C <program> <command>
			args := commands[i][1:]
			for j := 0; j+1 < len(args); j++ {
				if args[j] == "-C" && j+2 < len(args) {
					RegisterCompletion(args[j+2], args[j+1])
					j += 2
				} else if args[j] == "-p" {
					prog, ok := completionMap[args[j+1]]
					if ok {
						out += "complete -C '" + prog + "' " + args[j+1] + "\n"
					} else {
						out += "complete: " + args[j+1] + ": no completion specification\n"
					}
				} else if args[j] == "-r" {
					delete(completionMap, args[j+1])
				}
			}
			continue
		}

		// declare
		if command == "declare" {
			args := commands[i][1:]
			if strings.Contains(args[0], "=") {
				splitted := strings.Split(args[0], "=")
				key, val := splitted[0], splitted[1]
				addVar(key, val)
			}
			if len(args) >= 2 {
				if args[0] == "-p" {
					out += printVar(args[1])
				}
			}
			continue
		}
		// Run external command
		if ok, _ := Executable(command); ok {
			if !external_command(commands[i], os.Stdin, os.Stdout, os.Stderr) {
				success = false
			}
			continue
		}
		// Not found
		out += command + ": command not found\n"
		success = false
	}
	out += strings.Join(showJobs(true), "")
	return out, success
}
