package lexer

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/gogim1/goscript/file"
)

type Token struct {
	Location file.SourceLocation
	Source   file.Source
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
		src := file.NewSource("")
		currChar := l.source[l.currIndex]
		if unicode.IsDigit(currChar) || strings.ContainsRune("-+", currChar) {
			for l.currIndex < len(l.source) {
				currChar = l.source[l.currIndex]
				if unicode.IsDigit(currChar) || strings.ContainsRune("-+./", currChar) {
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
		} else if unicode.IsLetter(currChar) {
			for l.currIndex < len(l.source) {
				currChar = l.source[l.currIndex]
				if unicode.IsLetter(currChar) || unicode.IsDigit(currChar) || currChar == '_' {
					src = append(src, currChar)
					l.currLocation.Update(currChar)
					l.currIndex++
				} else {
					break
				}
			}
		} else if strings.ContainsRune("(){}[]=@&", currChar) {
			src = append(src, currChar)
			l.currLocation.Update(currChar)
			l.currIndex++
		} else if currChar == '"' {
			src = append(src, currChar)
			l.currLocation.Update(currChar)
			l.currIndex++
			for l.currIndex < len(l.source) {
				currChar = l.source[l.currIndex]
				if currChar != '"' || (currChar == '"' && countTrailingEscape(src)%2 != 0) {
					src = append(src, currChar)
					l.currLocation.Update(currChar)
					l.currIndex++
				} else {
					break
				}
			}
			if l.currIndex < len(l.source) && currChar == '"' {
				src = append(src, currChar)
				l.currLocation.Update(currChar)
				l.currIndex++
			} else {
				return Token{}, &file.Error{Location: tokenLocation, Message: "incomplete string literal"}
			}
		} else if currChar == '#' {
			for l.currIndex < len(l.source) {
				currChar = l.source[l.currIndex]
				if currChar != '\n' {
					l.currLocation.Update(currChar)
					l.currIndex++
				} else {
					break
				}
			}
			return l.nextToken()
		} else {
			return Token{}, &file.Error{Location: l.currLocation, Message: "unsupported token starting character"}
		}
		return Token{Location: tokenLocation, Source: src}, nil
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
		if reflect.DeepEqual(token, eof) {
			break
		} else {
			tokens = append(tokens, token)
		}
	}
	return tokens, nil
}
