package lexer_test

import (
	"strconv"
	"testing"

	"github.com/gogim1/goscript/file"
	. "github.com/gogim1/goscript/lexer"
	"github.com/stretchr/testify/assert"
)

func TestLex(t *testing.T) {
	tests := []struct {
		input  string
		tokens []Token
	}{
		{
			"",
			[]Token{},
		},
		{
			`1`,
			[]Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Source: "1"},
			},
		},
		{
			"+0.1 -1.1\r0/1\t-1/2\n",
			[]Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Source: "+0.1"},
				{Location: file.SourceLocation{Line: 1, Col: 6}, Source: "-1.1"},
				{Location: file.SourceLocation{Line: 1, Col: 11}, Source: "0/1"},
				{Location: file.SourceLocation{Line: 1, Col: 15}, Source: "-1/2"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tokens, err := Lex(file.NewSource(test.input))
			if err != nil {
				t.Errorf("%s:\n%v", strconv.Quote(test.input), err)
				return
			}
			if !compareTokens(tokens, test.tokens) {
				t.Errorf("%s:\ngot\n\t%+v\nexpected\n\t%+v", test.input, tokens, test.tokens)
			}
		})
	}
}

func compareTokens(i1, i2 []Token) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k] != i2[k] {
			return false
		}
	}
	return true
}

func TestLex_error(t *testing.T) {
	tests := []string{
		"♥︎",
		"1/0",
		"1.1.1",
	}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			_, err := Lex(file.NewSource(test))
			assert.NotNil(t, err)
		})
	}

}
