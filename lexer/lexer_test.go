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
		tokens []Token
	}{
		{
			"",
			[]Token{},
		},
		{
			`1`,
			[]Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Source: file.NewSource("1")},
			},
		},
		{
			"+0.1 -1.1\r0/1\t-1/2\n",
			[]Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Source: file.NewSource("+0.1")},
				{Location: file.SourceLocation{Line: 1, Col: 6}, Source: file.NewSource("-1.1")},
				{Location: file.SourceLocation{Line: 1, Col: 11}, Source: file.NewSource("0/1")},
				{Location: file.SourceLocation{Line: 1, Col: 15}, Source: file.NewSource("-1/2")},
			},
		},
		{
			"add",
			[]Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Source: file.NewSource("add")},
			},
		},
		{
			"int32 variable_name",
			[]Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Source: file.NewSource("int32")},
				{Location: file.SourceLocation{Line: 1, Col: 7}, Source: file.NewSource("variable_name")},
			},
		},
		{
			"(){}[]=@&",
			[]Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Source: file.NewSource("(")},
				{Location: file.SourceLocation{Line: 1, Col: 2}, Source: file.NewSource(")")},
				{Location: file.SourceLocation{Line: 1, Col: 3}, Source: file.NewSource("{")},
				{Location: file.SourceLocation{Line: 1, Col: 4}, Source: file.NewSource("}")},
				{Location: file.SourceLocation{Line: 1, Col: 5}, Source: file.NewSource("[")},
				{Location: file.SourceLocation{Line: 1, Col: 6}, Source: file.NewSource("]")},
				{Location: file.SourceLocation{Line: 1, Col: 7}, Source: file.NewSource("=")},
				{Location: file.SourceLocation{Line: 1, Col: 8}, Source: file.NewSource("@")},
				{Location: file.SourceLocation{Line: 1, Col: 9}, Source: file.NewSource("&")},
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
				{Location: file.SourceLocation{Line: 2, Col: 2}, Source: file.NewSource("1")},
			},
		},
		{
			`"s t r i n g" "\" \\\"" "#" "" "\\"`,
			[]Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Source: file.NewSource(`"s t r i n g"`)},
				{Location: file.SourceLocation{Line: 1, Col: 15}, Source: file.NewSource(`"\" \\\""`)},
				{Location: file.SourceLocation{Line: 1, Col: 25}, Source: file.NewSource(`"#"`)},
				{Location: file.SourceLocation{Line: 1, Col: 29}, Source: file.NewSource(`""`)},
				{Location: file.SourceLocation{Line: 1, Col: 32}, Source: file.NewSource(`"\\"`)},
			},
		},
		{
			`letrec ( sum = (add 100/11 61 +15/7 1.355 -41.06) ) { (put sum "\n") }`,
			[]Token{
				{Location: file.SourceLocation{Line: 1, Col: 1}, Source: file.NewSource(`letrec`)},
				{Location: file.SourceLocation{Line: 1, Col: 8}, Source: file.NewSource(`(`)},
				{Location: file.SourceLocation{Line: 1, Col: 10}, Source: file.NewSource(`sum`)},
				{Location: file.SourceLocation{Line: 1, Col: 14}, Source: file.NewSource(`=`)},
				{Location: file.SourceLocation{Line: 1, Col: 16}, Source: file.NewSource(`(`)},
				{Location: file.SourceLocation{Line: 1, Col: 17}, Source: file.NewSource(`add`)},
				{Location: file.SourceLocation{Line: 1, Col: 21}, Source: file.NewSource(`100/11`)},
				{Location: file.SourceLocation{Line: 1, Col: 28}, Source: file.NewSource(`61`)},
				{Location: file.SourceLocation{Line: 1, Col: 31}, Source: file.NewSource(`+15/7`)},
				{Location: file.SourceLocation{Line: 1, Col: 37}, Source: file.NewSource(`1.355`)},
				{Location: file.SourceLocation{Line: 1, Col: 43}, Source: file.NewSource(`-41.06`)},
				{Location: file.SourceLocation{Line: 1, Col: 49}, Source: file.NewSource(`)`)},
				{Location: file.SourceLocation{Line: 1, Col: 51}, Source: file.NewSource(`)`)},
				{Location: file.SourceLocation{Line: 1, Col: 53}, Source: file.NewSource(`{`)},
				{Location: file.SourceLocation{Line: 1, Col: 55}, Source: file.NewSource(`(`)},
				{Location: file.SourceLocation{Line: 1, Col: 56}, Source: file.NewSource(`put`)},
				{Location: file.SourceLocation{Line: 1, Col: 60}, Source: file.NewSource(`sum`)},
				{Location: file.SourceLocation{Line: 1, Col: 64}, Source: file.NewSource(`"\n"`)},
				{Location: file.SourceLocation{Line: 1, Col: 68}, Source: file.NewSource(`)`)},
				{Location: file.SourceLocation{Line: 1, Col: 70}, Source: file.NewSource(`}`)},
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
