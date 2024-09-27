package lexer

import "regexp"

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()-_=+[{]}\\|;:'\",<.>/? \t\n\r"

var numberRegexp = regexp.MustCompile(`^[+-]?((0|[1-9][0-9]*)|((0|[1-9][0-9]*)\.([0-9]*[1-9]))|((0|[1-9][0-9]*)/([1-9][0-9]*)))$`)
