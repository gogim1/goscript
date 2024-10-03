package lexer

import (
	"strings"
	"unicode"

	"github.com/gogim1/goscript/file"
)

type lexer struct {
	source       file.Source
	currIndex    int
	currLocation file.SourceLocation
}

func (l *lexer) nextToken() (*Token, *file.Error) {
	for l.currIndex < len(l.source) && unicode.IsSpace(l.source[l.currIndex]) {
		l.currLocation.Update(l.source[l.currIndex])
		l.currIndex++
	}
	begin := l.currIndex
	kind := Unknown
	if l.currIndex < len(l.source) {
		tokenLocation := l.currLocation
		currChar := l.source[l.currIndex]
		if unicode.IsDigit(currChar) || strings.ContainsRune("-+", currChar) {
			kind = Number
			for l.currIndex < len(l.source) {
				currChar = l.source[l.currIndex]
				if unicode.IsDigit(currChar) || strings.ContainsRune("-+./", currChar) {
					l.currLocation.Update(currChar)
					l.currIndex++
				} else {
					break
				}
			}
			if !numberRegexp.MatchString(string(l.source[begin:l.currIndex])) {
				return nil, &file.Error{Location: tokenLocation, Message: "invalid number literal"}
			}
		} else if unicode.IsLetter(currChar) {
			for l.currIndex < len(l.source) {
				currChar = l.source[l.currIndex]
				if unicode.IsLetter(currChar) || unicode.IsDigit(currChar) || currChar == '_' {
					l.currLocation.Update(currChar)
					l.currIndex++
				} else {
					break
				}
			}
			kind = Identifier
			for _, kw := range keyword {
				if string(l.source[begin:l.currIndex]) == kw {
					kind = Keyword
					break
				}
			}
		} else if strings.ContainsRune("(){}[]=@&", currChar) {
			kind = Symbol
			l.currLocation.Update(currChar)
			l.currIndex++
		} else if currChar == '"' {
			kind = String
			l.currLocation.Update(currChar)
			l.currIndex++
			for l.currIndex < len(l.source) {
				currChar = l.source[l.currIndex]
				if currChar != '"' || (currChar == '"' && countTrailingEscape(string(l.source[begin:l.currIndex]))%2 != 0) {
					l.currLocation.Update(currChar)
					l.currIndex++
				} else {
					break
				}
			}
			if l.currIndex < len(l.source) && currChar == '"' {
				l.currLocation.Update(currChar)
				l.currIndex++
			} else {
				return nil, &file.Error{Location: tokenLocation, Message: "incomplete string literal"}
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
			return nil, &file.Error{Location: l.currLocation, Message: "unsupported token starting character"}
		}
		return &Token{
			Location: tokenLocation,
			Kind:     kind,
			Source:   string(l.source[begin:l.currIndex]),
		}, nil
	} else {
		return &eof, nil
	}
}

func Lex(source file.Source) ([]*Token, *file.Error) {
	sl := file.SourceLocation{Line: 1, Col: 1}
	for _, char := range source {
		if !strings.ContainsRune(charSet, char) {
			return nil, &file.Error{Location: sl, Message: "unsupported character"}
		}
		sl.Update(char)
	}

	l := &lexer{
		source:       source,
		currLocation: file.SourceLocation{Line: 1, Col: 1},
	}
	tokens := make([]*Token, 0)
	for {
		token, err := l.nextToken()
		if err != nil {
			return nil, err
		}
		if token == &eof {
			break
		} else {
			tokens = append(tokens, token)
		}
	}
	return tokens, nil
}
