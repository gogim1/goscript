package runtime_test

import (
	"testing"

	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
	"github.com/gogim1/goscript/parser"
	. "github.com/gogim1/goscript/runtime"
	"github.com/stretchr/testify/assert"
)

func TestRuntime(t *testing.T) {
	tests := []struct {
		input, value string
	}{
		{`10/5`, `2`},
		{`letrec () {1.1}`, `11/10`},
		{`letrec (a=1 b=2) {"hello world"}`, `hello world`},
		{`letrec (a=1 b=2) {a}`, `1`},
		{`letrec (a=1 b="2") { letrec() {b} }`, `2`},
		{`letrec (a=1 b="2") { letrec(a=3) {a} }`, `3`},
		{`if 1 then 2 else 3`, `2`},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tokens, err := lexer.Lex(file.NewSource(test.input))
			assert.Nil(t, err)

			node, err := parser.Parse(tokens)
			assert.Nil(t, err)

			state := NewState(node)
			assert.Nil(t, state.Execute())
			assert.Equal(t, test.value, state.Value())
		})
	}
}

func TestRuntime_error(t *testing.T) {
	tests := []string{
		`if "true" then 2 else 3`,
	}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			tokens, err := lexer.Lex(file.NewSource(test))
			assert.Nil(t, err)

			node, err := parser.Parse(tokens)
			assert.Nil(t, err)

			state := NewState(node)
			assert.NotNil(t, state.Execute())
		})
	}
}
