package main

import (
	"errors"
	"fmt"
)

var shellVariables = map[string]string{}

func addVar(key, val string) error {
	if (key[0] >= 'A' && key[0] <= 'Z') || (key[0] >= 'a' && key[0] <= 'z') || key[0] == '_' {
		return errors.New("not a valid identifier")
	}
	shellVariables[key] = val
	return nil
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
