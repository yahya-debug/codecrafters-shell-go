package main

import (
	"bufio"
	"os"
)

var history []string
var histIndex int = -1
var l_append int = 0

// Transitions using arrow keys
func historyPrev() string {
	if len(history) == 0 || histIndex == -1 {
		return ""
	}
	if histIndex > 0 {
		histIndex--
	}
	return history[histIndex]
}
func historyNext() string {
	if len(history) == 0 {
		return ""
	}
	if histIndex < len(history)-1 {
		histIndex++
	} else {
		histIndex++
		return ""
	}
	return history[histIndex]
}

// Read history from a file
func ReadHist(file_name string) {
	file, err := os.OpenFile(file_name, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	file_reader := bufio.NewScanner(file)
	for file_reader.Scan() {
		line := file_reader.Text()
		history = append(history, line)
	}
	histIndex = len(history)
	l_append = len(history)
}

// Write History
func WriteHist(file_name string) {
	file, err := os.OpenFile(file_name, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	file_writer := bufio.NewWriter(file)
	for _, com := range history {
		file_writer.WriteString(com + "\n")
	}
	file_writer.Flush()
}
