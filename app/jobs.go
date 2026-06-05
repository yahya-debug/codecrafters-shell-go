package main

import (
	"fmt"
	"os"
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

func jobRun(inp []string) {
	if len(inp) == 0 {
		return
	}
	runInst := exec.Command(inp[0], inp[1:]...)
	runInst.Stdout = os.Stdout
	runInst.Stderr = os.Stderr
	err := runInst.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: command not found\n", inp[0])
		return
	}
	fmt.Printf("[%d] %d\n", len(jobs)+1, runInst.Process.Pid)

	os.Stdout.Sync()

	newJob := &Job{inp, runInst.Process.Pid, runInst, false}
	_ = addJob(newJob)

	go func(J *Job) { // run in background, line by line block by block then changes the status of completion
		_ = J.process.Wait()
		J.completed = true
	}(newJob)
}

func addJob(J *Job) int {
	jobs = append(jobs, J)
	return len(jobs)
}

func showJobs() []string {
	var all []string
	var line strings.Builder
	sz := len(jobs) - 1
	for i := sz; i >= 0; i-- {
		symb := "  "
		stat := "Running"
		r := i + 1

		if i == sz {
			symb = "+  "
		} else if i == sz-1 {
			symb = "-  "
		}

		if jobs[i].completed {
			stat = "Done"
		}

		var commandStr string
		if stat != "Running" {
			commandStr = strings.Join(jobs[i].command, " ")
		} else {
			commandStr = strings.Join(jobs[i].command, " ") + " &"
		}

		line.WriteString("[" + strconv.Itoa(r) + "]" + symb + stat + "\t\t" + commandStr + "\n")
		all = append(all, line.String())
		line.Reset()

		if stat == "Done" {
			jobs = append(jobs[:i], jobs[i+1:]...)
		}
	}

	return all
}
