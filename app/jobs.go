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
	runInst := exec.Command("bash", "-c", strings.Join(inp, " "))
	runInst.Stdin = os.Stdin
	runInst.Stdout = os.Stdout
	runInst.Stderr = os.Stderr
	err := runInst.Start()
	if err != nil {
		return
	}
	fmt.Print("\r[" + strconv.Itoa(len(jobs)+1) + "] " + strconv.Itoa(runInst.Process.Pid) + "\n")

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
