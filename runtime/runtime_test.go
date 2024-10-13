package runtime_test

import (
	"testing"

	"github.com/gogim1/goscript/ast"
	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
	"github.com/gogim1/goscript/parser"
	"github.com/gogim1/goscript/runtime"
	. "github.com/gogim1/goscript/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func lexAndParse(t *testing.T, src string) ast.ExprNode {
	tokens, err := lexer.Lex(file.NewSource(src))
	require.Nil(t, err)

	node, err := parser.Parse(tokens)
	require.Nil(t, err)

	return node
}

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
		{`(callcc lambda (k) { k })`, `<continuation evaluated at (SourceLocation 1 1)>`},
		{`(callcc lambda (k) { 1 })`, `1`},
		{`(callcc lambda (k) { [(k 1) 2] })`, `1`},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			state := NewState(lexAndParse(t, test.input))
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
			state := NewState(lexAndParse(t, test))
			assert.NotNil(t, state.Execute())
			assert.Equal(t, NewVoid(), state.Value())
		})
	}
}

func TestRuntimeInteraction(t *testing.T) {
	t.Run("call script function", func(t *testing.T) {
		src := `
		letrec (
			v = 1
		) {
			[
				(reg "test0" lambda(){ v })
				(reg "test1" lambda(v){ v })
				(reg "test1" lambda(v){ (add v 1) })  # doesn't work
				(reg "test2" lambda(v){ (add v 1) })  # v should be Number
			]
		}
			`

		state := runtime.NewState(lexAndParse(t, src))
		err := state.Execute()
		require.Nil(t, err)

		v, err := state.Call("test0")
		assert.Nil(t, err)
		assert.True(t, v != nil && v.String() == `1`)

		v, err = state.Call("test1", 42)
		assert.Nil(t, err)
		assert.True(t, v != nil && v.String() == `42`)

		v, err = state.Call("test0", 1)
		assert.NotNil(t, err)
		assert.Nil(t, v)

		v, err = state.Call("test3")
		assert.NotNil(t, err)
		assert.Nil(t, v)

		v, err = state.Call("test2", "string")
		assert.NotNil(t, err)
		assert.Nil(t, v)
	})

	t.Run("call golang function", func(t *testing.T) {
		plus1 := func(args ...runtime.Value) runtime.Value {
			arg := args[0].(*Number)
			return runtime.NewNumber(arg.Numerator+arg.Denominator, arg.Numerator)
		}

		src := `(go "plus1" 1)`
		state := runtime.NewState(lexAndParse(t, src)).Register("plus1", plus1)
		assert.Nil(t, state.Execute())
		assert.Equal(t, `2`, state.Value().String())

		plus1 = func(args ...runtime.Value) runtime.Value {
			return runtime.NewNumber(0, 1)
		}
		assert.Nil(t, state.Execute())
		assert.Equal(t, `2`, state.Value().String())

		str := func(args ...runtime.Value) runtime.Value {
			return runtime.NewString("string")
		}
		src = `(go "str")`
		state = runtime.NewState(lexAndParse(t, src)).Register("str", str)
		assert.Nil(t, state.Execute())
		assert.Equal(t, `string`, state.Value().String())

		state = runtime.NewState(lexAndParse(t, src)).Register("str", plus1).Register("str", str)
		assert.Nil(t, state.Execute())
		assert.Equal(t, `string`, state.Value().String())

		src = `(go "str" (go "plus1" 1))`
		state = runtime.NewState(lexAndParse(t, src)).Register("plus1", plus1).Register("str", str)
		assert.Nil(t, state.Execute())
		assert.Equal(t, `string`, state.Value().String())

		src = `(go "unknown")`
		state = runtime.NewState(lexAndParse(t, src))
		assert.NotNil(t, state.Execute())
		assert.Equal(t, NewVoid(), state.Value())

		// TODO: should we handle panics?
		// raise := func(args ...runtime.Value) runtime.Value {
		// 	panic("error")
		// }
		// src = `(go "raise")`
		// state = runtime.NewState(lexAndParse(t, src)).Register("raise", raise)
		// assert.NotNil(t, state.Execute())
	})
}
