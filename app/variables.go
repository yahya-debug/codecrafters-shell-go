package main

import "fmt"

var shellVariables = map[string]string{}

func addVar(key, val string) {
	shellVariables[key] = val
}

func removeVar(key string) {
	delete(shellVariables, key)
}

func printVar(key string) string {
	val, ok := shellVariables[key]
	if ok {
		return fmt.Sprintf("declare -- %s=\"%s\"\n", key, val)
	} else {
		return fmt.Sprintf("declare: %s: not found\n", key)
	}
}
