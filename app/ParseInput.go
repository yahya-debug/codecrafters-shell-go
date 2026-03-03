package main

import (
	"strings"
)

// Parse input to split the arguments using '"' '\” '>'
func ParseInput(str string) []string {
	var res []string
	var cur strings.Builder
	var inSingle, inDouble bool = false, false
	for i := 0; i < len(str); i++ {
		if i < len(str)-1 && str[i] == '\\' {
			cur.WriteByte(str[i+1])
			i++
			continue
		}
		switch str[i] {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
				continue
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
				continue
			}
		case ' ':
			if !inSingle && !inDouble {
				if cur.Len() > 0 {
					res = append(res, cur.String())
					cur.Reset()
				}
				continue
			}
		case '>':
			if !inSingle && !inDouble {
				res = append(res, cur.String())
				cur.Reset()
			}
			res = append(res, ">")
			continue

		}
		cur.WriteByte(str[i])
	}
	if cur.Len() > 0 {
		res = append(res, cur.String())
	}
	return res
}
