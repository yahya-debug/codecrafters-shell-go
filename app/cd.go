package main

import "os"

func HandleCD(path string) bool {
	err := os.Chdir(path)
	if err != nil {
		return false
	}
	return true
}
