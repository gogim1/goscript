package file

import "fmt"

type Source []rune

func NewSource(contents string) Source {
	return []rune(contents)
}

func (s Source) String() string {
	return string(s)
}

type SourceLocation struct {
	Line int
	Col  int
}

func (sl SourceLocation) String() string {
	if sl.Line <= 0 || sl.Col <= 0 {
		return "(SourceLocation N/A)"
	}
	return fmt.Sprintf("(SourceLocation %d %d)", sl.Line, sl.Col)
}

func (sl *SourceLocation) Update(char rune) {
	if char == '\n' {
		sl.Line = sl.Line + 1
		sl.Col = 1
	} else {
		sl.Col = sl.Col + 1
	}
}
