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
		{
			"add",
			[]Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Source: "add"},
			},
		},
		{
			"int32 variable_name",
			[]Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Source: "int32"},
				{Location: file.SourceLocation{Line: 1, Col: 7}, Source: "variable_name"},
			},
		},
		{
			"(){}[]=@&",
			[]Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Source: "("},
				{Location: file.SourceLocation{Line: 1, Col: 2}, Source: ")"},
				{Location: file.SourceLocation{Line: 1, Col: 3}, Source: "{"},
				{Location: file.SourceLocation{Line: 1, Col: 4}, Source: "}"},
				{Location: file.SourceLocation{Line: 1, Col: 5}, Source: "["},
				{Location: file.SourceLocation{Line: 1, Col: 6}, Source: "]"},
				{Location: file.SourceLocation{Line: 1, Col: 7}, Source: "="},
				{Location: file.SourceLocation{Line: 1, Col: 8}, Source: "@"},
				{Location: file.SourceLocation{Line: 1, Col: 9}, Source: "&"},
			},
		},
		{
			"#comment",
			[]Token{},
		},
		{
			"# #comment \n #\n",
			[]Token{},
		},
		{
			"# #comment \n 1 #\n",
			[]Token{
				{Location: file.SourceLocation{Line: 2, Col: 2}, Source: "1"},
			},
		},
		{
			`"s t r i n g" "\" \\\"" "#" "" "\\"`,
			[]Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Source: `"s t r i n g"`},
				{Location: file.SourceLocation{Line: 1, Col: 15}, Source: `"\" \\\""`},
				{Location: file.SourceLocation{Line: 1, Col: 25}, Source: `"#"`},
				{Location: file.SourceLocation{Line: 1, Col: 29}, Source: `""`},
				{Location: file.SourceLocation{Line: 1, Col: 32}, Source: `"\\"`},
			},
		},
		{
			`letrec ( sum = (add 100/11 61 +15/7 1.355 -41.06) ) { (put sum "\n") }`,
			[]Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Source: `letrec`},
				{Location: file.SourceLocation{Line: 1, Col: 8}, Source: `(`},
				{Location: file.SourceLocation{Line: 1, Col: 10}, Source: `sum`},
				{Location: file.SourceLocation{Line: 1, Col: 14}, Source: `=`},
				{Location: file.SourceLocation{Line: 1, Col: 16}, Source: `(`},
				{Location: file.SourceLocation{Line: 1, Col: 17}, Source: `add`},
				{Location: file.SourceLocation{Line: 1, Col: 21}, Source: `100/11`},
				{Location: file.SourceLocation{Line: 1, Col: 28}, Source: `61`},
				{Location: file.SourceLocation{Line: 1, Col: 31}, Source: `+15/7`},
				{Location: file.SourceLocation{Line: 1, Col: 37}, Source: `1.355`},
				{Location: file.SourceLocation{Line: 1, Col: 43}, Source: `-41.06`},
				{Location: file.SourceLocation{Line: 1, Col: 49}, Source: `)`},
				{Location: file.SourceLocation{Line: 1, Col: 51}, Source: `)`},
				{Location: file.SourceLocation{Line: 1, Col: 53}, Source: `{`},
				{Location: file.SourceLocation{Line: 1, Col: 55}, Source: `(`},
				{Location: file.SourceLocation{Line: 1, Col: 56}, Source: `put`},
				{Location: file.SourceLocation{Line: 1, Col: 60}, Source: `sum`},
				{Location: file.SourceLocation{Line: 1, Col: 64}, Source: `"\n"`},
				{Location: file.SourceLocation{Line: 1, Col: 68}, Source: `)`},
				{Location: file.SourceLocation{Line: 1, Col: 70}, Source: `}`},
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
		`"`,
		`"\"`,
		`"\\\"`,
	}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			_, err := Lex(file.NewSource(test))
			assert.NotNil(t, err)
		})
	}

}
