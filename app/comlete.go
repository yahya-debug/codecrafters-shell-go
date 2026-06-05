package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// completionMap maps a command name to the program that generates its completions.
var completionMap = map[string]string{}

func RegisterCompletion(command, program string) {
	completionMap[command] = program
}

// runCompletionProgram calls the -C program and returns one candidate per output line.
// Env: COMP_LINE, COMP_POINT. Args: program word prev-word.
func runCompletionProgram(program, line string, point int, word, prevWord string) []string {
	cmd := exec.Command(program, word, prevWord)
	cmd.Env = append(os.Environ(),
		"COMP_LINE="+line,
		"COMP_POINT="+strconv.Itoa(point),
	)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	var results []string
	for _, s := range strings.Split(string(out), "\n") {
		if s != "" {
			results = append(results, s)
		}
	}
	return results
}

// completionCandidates returns -C candidates for the current line/point,
// or nil if no -C mapping exists for the leading command.
func completionCandidates(line string, point int) []string {
	if point == 0 {
		return nil
	}
	parts := strings.Fields(line[:point])
	if len(parts) == 0 {
		return nil
	}
	prog, ok := completionMap[parts[0]]
	if !ok {
		return nil
	}

	word := ""
	prevWord := parts[0]
	if len(parts) >= 2 {
		if line[point-1] != ' ' {
			word = parts[len(parts)-1]
			if len(parts) >= 3 {
				prevWord = parts[len(parts)-2]
			}
		} else {
			prevWord = parts[len(parts)-1]
		}
	}

	return runCompletionProgram(prog, line[:point], point, word, prevWord)
}
