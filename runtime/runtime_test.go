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
		{`"hello world"`, `hello world`},
		{`[1/1 "2" 3]`, `3`},
		{`if 1 then 2 else 3`, `2`},
		{`letrec () {1.1}`, `11/10`},
		{`letrec (a=1 b=2) {"hello world"}`, `hello world`},
		{`letrec (a=1 b=2) {a}`, `1`},
		{`letrec (a=1 b="2") { letrec() {b} }`, `2`},
		{`letrec (a=1 b="2") { letrec(a=3) {a} }`, `3`},
		{`letrec (a=1 b="2") { [letrec(a=3) {a} a] }`, `1`},
		{`letrec (a=1 b=lambda() {c} c=2 ) { (b) }`, `2`},
		{`letrec (A=2 c=lambda() {A}) { (c) }`, `2`},
		{`lambda () {1}`, `<closure evaluated at (SourceLocation 1 1)>`},
		{`(lambda () {1})`, `1`},
		{`(lambda (a b) {a} 1 2)`, `1`},
		{`((lambda (a) { lambda (a) { a } } 1) 2)`, `2`},
		{`(lambda (a) { (lambda (b c) { c } a a) } 1)`, `1`},
		{`(lambda (a) { [(lambda (a) { a } 1) a] } 2)`, `2`},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tokens, err := lexer.Lex(file.NewSource(test.input))
			assert.Nil(t, err)

			node, err := parser.Parse(tokens)
			assert.Nil(t, err)

			state := NewState(node)
			assert.Nil(t, state.Execute())
			assert.Equal(t, test.value, state.Value().String())
		})
	}
}

func TestRuntime_error(t *testing.T) {
	tests := []string{
		`if "true" then 2 else 3`,
		`[letrec (a = 1) {a} a]`,
		`(lambda () { 1 } 2)`,
		`letrec (a = 1) {(a)}`,
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
