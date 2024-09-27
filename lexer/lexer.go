package lexer

import (
	"strings"
	"unicode"

	"github.com/gogim1/goscript/file"
)

type Token struct {
	Location file.SourceLocation
	Source   string
}

type lexer struct {
	source       file.Source
	currIndex    int
	currLocation file.SourceLocation
}

var (
	eof = Token{}
)

func (l *lexer) nextToken() (Token, *file.Error) {
	for l.currIndex < len(l.source) && unicode.IsSpace(l.source[l.currIndex]) {
		l.currLocation.Update(l.source[l.currIndex])
		l.currIndex++
	}
	if l.currIndex < len(l.source) {
		tokenLocation := l.currLocation
		src := []rune("")
		currChar := l.source[l.currIndex]
		if unicode.IsDigit(currChar) || currChar == '-' || currChar == '+' {
			for l.currIndex < len(l.source) {
				currChar = l.source[l.currIndex]
				if unicode.IsDigit(currChar) || currChar == '-' || currChar == '+' || currChar == '.' || currChar == '/' {
					src = append(src, currChar)
					l.currLocation.Update(currChar)
					l.currIndex++
				} else {
					break
				}
			}
			if !numberRegexp.MatchString(string(src)) {
				return Token{}, &file.Error{Location: tokenLocation, Message: "invalid number literal"}
			}
		} else {
			return Token{}, &file.Error{Location: l.currLocation, Message: "unsupported token starting character"}
		}
		return Token{Location: tokenLocation, Source: string(src)}, nil
	} else {
		return eof, nil
	}
}

func Lex(source file.Source) ([]Token, *file.Error) {
	sl := file.SourceLocation{Line: 1, Col: 1}
	for _, char := range source {
		if !strings.ContainsRune(charSet, char) {
			return nil, &file.Error{Location: sl, Message: "unsupported character"}
		}
		sl.Update(char)
	}

	l := &lexer{
		source:       source,
		currIndex:    0,
		currLocation: file.SourceLocation{Line: 1, Col: 1},
	}
	tokens := make([]Token, 0)
	for {
		token, err := l.nextToken()
		if err != nil {
			return nil, err
		}
		if token == eof {
			break
		} else {
			tokens = append(tokens, token)
		}
	}
	return tokens, nil
}
