package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var jobs []*Job

type Job struct {
	order   int
	command []string
	PID     int
	process *exec.Cmd
	doneCh  chan struct{}
}

func (j *Job) isDone() bool {
	select {
	case <-j.doneCh:
		return true
	default:
		return false
	}
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
		return
	}

	var ord int
	if len(jobs) == 0 {
		ord = 1
	} else {
		ord = jobs[len(jobs)-1].order + 1
	}

	fmt.Printf("[%d] %d\n", ord, runInst.Process.Pid)
	os.Stdout.Sync()

	newJob := &Job{ord, inp, runInst.Process.Pid, runInst, make(chan struct{})}
	_ = addJob(newJob)

	go func(J *Job) {
		_ = J.process.Wait()
		close(J.doneCh)
	}(newJob)
}

func addJob(J *Job) int {
	jobs = append(jobs, J)
	return len(jobs)
}

func showJobs(filterDone bool) []string {
	// Yield so any completion goroutines can run before we check status.
	runtime.Gosched()

	var all []string
	var line strings.Builder
	sz := len(jobs) - 1
	for i := sz; i >= 0; i-- {
		symb := "  "
		stat := "Running"

		if i == sz {
			symb = "+  "
		} else if i == sz-1 {
			symb = "-  "
		}

		if jobs[i].isDone() {
			stat = "Done"
		}

		var commandStr string
		if stat == "Running" {
			commandStr = strings.Join(jobs[i].command, " ") + " &"
		} else {
			commandStr = strings.Join(jobs[i].command, " ")
		}

		line.WriteString("[" + strconv.Itoa(jobs[i].order) + "]" + symb + stat + "\t\t" + commandStr + "\n")
		if !filterDone || stat == "Done" {
			all = append(all, line.String())
		}
		line.Reset()

		if stat == "Done" {
			jobs = append(jobs[:i], jobs[i+1:]...)
		}
	}

	return all
}
