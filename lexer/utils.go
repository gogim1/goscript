package lexer

import "regexp"

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()-_=+[{]}\\|;:'\",<.>/? \t\n\r"

var numberRegexp = regexp.MustCompile(`^[+-]?((0|[1-9][0-9]*)|((0|[1-9][0-9]*)\.([0-9]*[1-9]))|((0|[1-9][0-9]*)/([1-9][0-9]*)))$`)

func countTrailingEscape(s string) int {
	cnt := 0
	pos := len(s) - 1
	for pos >= 0 {
		if s[pos] == '\\' {
			cnt++
			pos--
		} else {
			break
		}
	}
	return cnt
}
