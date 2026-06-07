package main

import (
	"strings"
)

// Parse input to split the arguments using '"' '\” '>'
func ParseInput(str string) ([]string, bool) {
	var res []string
	var cur strings.Builder
	var inSingle, inDouble bool = false, false
	multi := false
	for i := 0; i < len(str); i++ {
		if !inSingle && i < len(str)-1 && str[i] == '\\' {
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
		case '|':
			if !inSingle && !inDouble {
				if cur.Len() > 0 {
					res = append(res, cur.String())
					cur.Reset()
				}
				res = append(res, "|")
			}
			continue
		case '>':
			if !inSingle && !inDouble {
				if cur.Len() > 0 {
					res = append(res, cur.String())
					cur.Reset()
				}
				if i < len(str)-1 && str[i+1] == '>' {
					res = append(res, ">>")
					i++
					continue
				}
			}
			res = append(res, ">")
			continue
		case '&':
			if !inSingle && !inDouble {
				if i < len(str)-1 && str[i+1] == '&' {
					res = append(res, "&&")
					multi = true
					i++
					continue
				}
			}
		case '1':
			if i < len(str)-1 && str[i+1] == '>' {
				if cur.Len() > 0 {
					res = append(res, cur.String())
					cur.Reset()
				}
				continue
			}
			cur.WriteByte(str[i])
			continue
		case '2':
			if i < len(str)-2 && str[i+1] == '>' && str[i+2] == '>' {
				if cur.Len() > 0 {
					res = append(res, cur.String())
					cur.Reset()
				}
				res = append(res, "2>>")
				i += 2
				continue
			}
			if i < len(str)-1 && str[i+1] == '>' {
				if cur.Len() > 0 {
					res = append(res, cur.String())
					cur.Reset()
				}
				res = append(res, "2>")
				i++
				continue
			}
			cur.WriteByte(str[i])
			continue
		}
		cur.WriteByte(str[i])
	}

	if cur.Len() > 0 {
		res = append(res, cur.String())
	}
	for it := 0; it < len(res); it++ {
		if strings.Contains(res[it], "${") {
			for i := 0; i < len(res[it]); i++ {
				if i+1 < len(res[it]) && res[it][i] == '$' && res[it][i+1] == '{' {
					var key strings.Builder
					var j int
					for j = i + 2; j < len(res[it]); j++ {
						if res[it][j] == '}' {
							break
						}
						key.WriteByte(res[it][j])
					}
					res[it] = res[it][:i] + shellVariables[key.String()] + res[it][j+1:]
					key.Reset()
				}
			}
		} else if strings.Contains(res[it], "$") {
			res[it] = strings.Split(res[it], "$")[0] + shellVariables[strings.Split(res[it], "$")[1]]
		}
		if len(strings.TrimSpace(res[it])) == 0 {
			if it+1 < len(res) {
				res = append(res[:it], res[it+1:]...)
				it--
			} else {
				res = res[:it]
			}
		}
	}

	return res, multi
}
