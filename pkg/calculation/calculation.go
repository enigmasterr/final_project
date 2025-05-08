package calculation

import (
	"strings"

	"github.com/enigmasterr/final_project/pkg/calculation"
)

func Calc(expression string) ([]string, error) {

	prior := map[string]int{
		"(": 0,
		")": 1,
		"+": 2,
		"-": 2,
		"*": 3,
		"/": 3,
	}
	var ans []string
	var st []string
	num := ""
	charset := "+-*/()0123456789"
	strange := false
	expression = strings.ReplaceAll(expression, " ", "")
	for _, sim := range expression {
		if !strings.ContainsRune(charset, sim) {
			strange = true
		}
	}
	if strange {
		return []string{}, calculation.ErrStrangeSymbols
	}
	for _, sim := range expression {
		if sim == '(' {
			if len(num) > 0 {
				ans = append(ans, num)
			}
			st = append(st, string(sim))
		} else {
			if sim == '+' || sim == '-' || sim == '*' || sim == '/' {
				if num != "" {
					ans = append(ans, num)
					num = ""
				}
				if len(st) == 0 {
					st = append(st, string(sim))
				} else {
					if prior[string(sim)] > prior[st[len(st)-1]] {
						st = append(st, string(sim))
					} else {
						for len(st) > 0 && prior[string(sim)] <= prior[st[len(st)-1]] {
							ans = append(ans, st[len(st)-1])
							st = st[:len(st)-1]
						}
						st = append(st, string(sim))
					}
				}
			} else if sim == ')' {
				if len(num) > 0 {
					ans = append(ans, num)
					num = ""
				}
				for st[len(st)-1] != "(" {
					ans = append(ans, st[len(st)-1])
					st = st[:len(st)-1]
				}
				st = st[:len(st)-1]
			} else {
				num += string(sim)
			}
		}
	}
	if num != "" {
		ans = append(ans, num)
		num = ""
	}
	for len(st) > 0 {
		if st[len(st)-1] == "(" || st[len(st)-1] == ")" {
			return []string{}, calculation.ErrInvalidExpression
		} else {
			ans = append(ans, st[len(st)-1])
			st = st[:len(st)-1]
		}
	}

	return ans, nil
}
