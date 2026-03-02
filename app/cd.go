package main

import (
	"os"
	"strings"
)

func HandleCD(path string) bool {
	if strings.ToLower(strings.TrimSpace(path)) == "~" {
		dir, _ := os.UserHomeDir()
		os.Chdir(dir)
		return true
	}
	err := os.Chdir(path)
	if err != nil {
		return false
	}
	return true
}
