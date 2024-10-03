package lexer_test

import (
	"testing"

	"github.com/gogim1/goscript/file"
	. "github.com/gogim1/goscript/lexer"
	"github.com/stretchr/testify/assert"
)

func TestLex(t *testing.T) {
	tests := []struct {
		input  string
		tokens []*Token
	}{
		{
			"",
			[]*Token{},
		},
		{
			`1`,
			[]*Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Kind: Number, Source: "1"},
			},
		},
		{
			"+0.1 -1.1\r0/1\t-1/2\n",
			[]*Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Kind: Number, Source: "+0.1"},
				{Location: file.SourceLocation{Line: 1, Col: 6}, Kind: Number, Source: "-1.1"},
				{Location: file.SourceLocation{Line: 1, Col: 11}, Kind: Number, Source: "0/1"},
				{Location: file.SourceLocation{Line: 1, Col: 15}, Kind: Number, Source: "-1/2"},
			},
		},
		{
			"add",
			[]*Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Kind: Identifier, Source: "add"},
			},
		},
		{
			"int32 variable_name",
			[]*Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Kind: Identifier, Source: "int32"},
				{Location: file.SourceLocation{Line: 1, Col: 7}, Kind: Identifier, Source: "variable_name"},
			},
		},
		{
			"(){}[]=@&",
			[]*Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Kind: Symbol, Source: "("},
				{Location: file.SourceLocation{Line: 1, Col: 2}, Kind: Symbol, Source: ")"},
				{Location: file.SourceLocation{Line: 1, Col: 3}, Kind: Symbol, Source: "{"},
				{Location: file.SourceLocation{Line: 1, Col: 4}, Kind: Symbol, Source: "}"},
				{Location: file.SourceLocation{Line: 1, Col: 5}, Kind: Symbol, Source: "["},
				{Location: file.SourceLocation{Line: 1, Col: 6}, Kind: Symbol, Source: "]"},
				{Location: file.SourceLocation{Line: 1, Col: 7}, Kind: Symbol, Source: "="},
				{Location: file.SourceLocation{Line: 1, Col: 8}, Kind: Symbol, Source: "@"},
				{Location: file.SourceLocation{Line: 1, Col: 9}, Kind: Symbol, Source: "&"},
			},
		},
		{
			"#comment",
			[]*Token{},
		},
		{
			"# #comment \n #\n",
			[]*Token{},
		},
		{
			"# #comment \n 1 #\n",
			[]*Token{
				{Location: file.SourceLocation{Line: 2, Col: 2}, Kind: Number, Source: "1"},
			},
		},
		{
			`"s t r i n g" "\" \\\"" "#" "" "\\"`,
			[]*Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Kind: String, Source: `"s t r i n g"`},
				{Location: file.SourceLocation{Line: 1, Col: 15}, Kind: String, Source: `"\" \\\""`},
				{Location: file.SourceLocation{Line: 1, Col: 25}, Kind: String, Source: `"#"`},
				{Location: file.SourceLocation{Line: 1, Col: 29}, Kind: String, Source: `""`},
				{Location: file.SourceLocation{Line: 1, Col: 32}, Kind: String, Source: `"\\"`},
			},
		},
		{
			`letrec ( sum = (add 100/11 61 +15/7 1.355 -41.06) ) { (put sum "\n") }`,
			[]*Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Kind: Keyword, Source: `letrec`},
				{Location: file.SourceLocation{Line: 1, Col: 8}, Kind: Symbol, Source: `(`},
				{Location: file.SourceLocation{Line: 1, Col: 10}, Kind: Identifier, Source: `sum`},
				{Location: file.SourceLocation{Line: 1, Col: 14}, Kind: Symbol, Source: `=`},
				{Location: file.SourceLocation{Line: 1, Col: 16}, Kind: Symbol, Source: `(`},
				{Location: file.SourceLocation{Line: 1, Col: 17}, Kind: Identifier, Source: `add`},
				{Location: file.SourceLocation{Line: 1, Col: 21}, Kind: Number, Source: `100/11`},
				{Location: file.SourceLocation{Line: 1, Col: 28}, Kind: Number, Source: `61`},
				{Location: file.SourceLocation{Line: 1, Col: 31}, Kind: Number, Source: `+15/7`},
				{Location: file.SourceLocation{Line: 1, Col: 37}, Kind: Number, Source: `1.355`},
				{Location: file.SourceLocation{Line: 1, Col: 43}, Kind: Number, Source: `-41.06`},
				{Location: file.SourceLocation{Line: 1, Col: 49}, Kind: Symbol, Source: `)`},
				{Location: file.SourceLocation{Line: 1, Col: 51}, Kind: Symbol, Source: `)`},
				{Location: file.SourceLocation{Line: 1, Col: 53}, Kind: Symbol, Source: `{`},
				{Location: file.SourceLocation{Line: 1, Col: 55}, Kind: Symbol, Source: `(`},
				{Location: file.SourceLocation{Line: 1, Col: 56}, Kind: Identifier, Source: `put`},
				{Location: file.SourceLocation{Line: 1, Col: 60}, Kind: Identifier, Source: `sum`},
				{Location: file.SourceLocation{Line: 1, Col: 64}, Kind: String, Source: `"\n"`},
				{Location: file.SourceLocation{Line: 1, Col: 68}, Kind: Symbol, Source: `)`},
				{Location: file.SourceLocation{Line: 1, Col: 70}, Kind: Symbol, Source: `}`},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tokens, err := Lex(file.NewSource(test.input))
			assert.Nil(t, err)

			assert.Equal(t, test.tokens, tokens)
		})
	}
}

func TestLex_error(t *testing.T) {
	tests := []string{
		"♥︎",
		"1/0",
		"1/-3",
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
