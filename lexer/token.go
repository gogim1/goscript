package lexer

import "github.com/gogim1/goscript/file"

type Kind int

const (
	Unknown Kind = iota
	Keyword
	Identifier
	Number
	String
	Symbol
)

type Token struct {
	Location file.SourceLocation
	Kind     Kind
	Source   string
}

var eof = Token{}

var keyword = []string{
	"if", "then", "else",
	"letrec",
	"lambda",
}
