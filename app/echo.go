package main

import (
	"fmt"
	"strings"
)

func HandleEcho(toPrint string) {
	var res []string
	var cur strings.Builder
	q := false
	for i := 0; i < len(toPrint); i++ {
		if toPrint[i] == '\'' {
			q = !q
			continue
		}
		if toPrint[i] == ' ' && !q {
			if cur.Len() > 0 {
				res = append(res, cur.String())
				cur.Reset()
			}
			continue
		}
		cur.WriteByte(toPrint[i])
	}
	if cur.Len() > 0 {
		res = append(res, cur.String())
	}
	fmt.Print(strings.Join(res, " "))
}
