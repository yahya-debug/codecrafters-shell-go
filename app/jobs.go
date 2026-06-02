package main

import (
	"os/exec"
	"strconv"
	"strings"
)

var jobs []*Job

type Job struct {
	command   []string
	PID       int
	process   *exec.Cmd
	completed bool
}

func addJob(J *Job) int {
	jobs = append(jobs, J)
	return len(jobs)
}

func showJobs() []string {
	var all []string
	var line strings.Builder
	for i := len(jobs) - 1; i >= 0; i-- {
		stat := "Running"
		r := i + 1

		if jobs[i].completed {
			stat = "Done"
		}

		line.WriteString("[" + strconv.Itoa(r) + "]  " + stat + "\t\t" + strings.Join(jobs[i].command, " ") + "\n")
		all = append(all, line.String())
		line.Reset()

		if stat == "Done" {
			jobs = append(jobs[:i], jobs[i+1:]...)
		}
	}

	return all
}
