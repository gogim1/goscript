package runtime_test

import (
	"testing"

	"github.com/gogim1/goscript/ast"
	"github.com/gogim1/goscript/conf"
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
		{`if 0 then 2 else 3`, `3`},
		{`letrec () {1.1}`, `11/10`},
		{`letrec (a=1 b=2) {"hello world"}`, `hello world`},
		{`letrec (a=1 b=2) {a}`, `1`},
		{`letrec (a=1 b="2") { letrec() {b} }`, `2`},
		{`letrec (a=1 b="2") { letrec(a=3) {a} }`, `3`},
		{`letrec (a=1 b="2") { [letrec(a=3) {a} a] }`, `1`},
		{`letrec (a=1 b=a) { b }`, `1`},
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
		{`&v letrec (v=1) {lambda () { 1 }}`, `1`},
		{`&v (lambda (v) { lambda () { 0 } } 1)`, `1`},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			state := NewState(lexAndParse(t, test.input), conf.New())
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
		`(void 1)`,
		`(isvoid)`,
		`(isnum)`,
		`(isstr)`,
		`(isclo)`,
		`(iscont)`,
		`(add 1)`,
		`(add "1" 2)`,
		`(sub "1" 2)`,
		`(mul 1 "2")`,
		`(div 1 "2")`,
		`(div 1 0)`,
		`(lt "1" "2")`,
		`(gt "1" "2")`,
		`(le "1" "2")`,
		`(ge "1" "2")`,
		`(eq 1 "2")`,
		`(eq "1" (void))`,
		`(ne 1 "2")`,
		`(ne "1" (void))`,
		`(and "1" "1")`,
		`(or "1" "1")`,
		`(not "1")`,
		`(not "1")`,
		`(getline "1")`,
		`(put)`,
		`(eval 1)`,
		`(callcc 1)`,
		`(reg "func" 1)`,
		`(go 1)`,
		`(concat 1 2)`,
		`&v lambda () { 1 }`,
		`&v "string"`,
	}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			state := NewState(lexAndParse(t, test), conf.New())
			err := state.Execute()
			assert.NotNil(t, err)
			assert.Equal(t, "<void>", state.Value().String())
			t.Log(err)
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

		state := runtime.NewState(lexAndParse(t, src), conf.New())
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
		conf := conf.New()
		plus1 := func(args ...runtime.Value) runtime.Value {
			arg := args[0].(*Number)
			return runtime.NewNumber(arg.Numerator+arg.Denominator, arg.Numerator)
		}

		src := `(go "plus1" 1)`
		state := runtime.NewState(lexAndParse(t, src), conf).Register("plus1", plus1)
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
		state = runtime.NewState(lexAndParse(t, src), conf).Register("str", str)
		assert.Nil(t, state.Execute())
		assert.Equal(t, `string`, state.Value().String())

		state = runtime.NewState(lexAndParse(t, src), conf).Register("str", plus1).Register("str", str)
		assert.Nil(t, state.Execute())
		assert.Equal(t, `string`, state.Value().String())

		src = `(go "str" (go "plus1" 1))`
		state = runtime.NewState(lexAndParse(t, src), conf).Register("plus1", plus1).Register("str", str)
		assert.Nil(t, state.Execute())
		assert.Equal(t, `string`, state.Value().String())

		src = `(go "unknown")`
		state = runtime.NewState(lexAndParse(t, src), conf)
		assert.NotNil(t, state.Execute())
		assert.Equal(t, "<void>", state.Value().String())

		// TODO: should we handle panics?
		// raise := func(args ...runtime.Value) runtime.Value {
		// 	panic("error")
		// }
		// src = `(go "raise")`
		// state = runtime.NewState(lexAndParse(t, src)).Register("raise", raise)
		// assert.NotNil(t, state.Execute())
	})
}

func TestIntrinsics(t *testing.T) {
	tests := []struct {
		input, value string
	}{
		{`(void)`, `<void>`},
		{`(isvoid 1)`, `0`},
		{`(isvoid (void))`, `1`},
		{`(isnum 0)`, `1`},
		{`(isnum "1")`, `0`},
		{`(isstr 1)`, `0`},
		{`(isstr "1")`, `1`},
		{`(isclo 1)`, `0`},
		{`(isclo lambda () { 0 })`, `1`},
		{`(iscont (callcc lambda(k) { (k k) }))`, `1`},
		{`(iscont (callcc lambda(k) { 1 }))`, `0`},
		{`(add 1 2)`, `3`},
		{`(add 0.3 2/3)`, `29/30`},
		{`(sub 1 2)`, `-1`},
		{`(mul 1 2)`, `2`},
		{`(div 9 -3)`, `-3`},
		{`(lt 1 2)`, `1`},
		{`(lt 2 1)`, `0`},
		{`(gt 1 2)`, `0`},
		{`(gt 2 1)`, `1`},
		{`(le 2 2)`, `1`},
		{`(le 2 1)`, `0`},
		{`(ge 1 2)`, `0`},
		{`(ge 1 1)`, `1`},
		{`(eq 1 1)`, `1`},
		{`(eq 1 2)`, `0`},
		{`(eq "1" "1")`, `1`},
		{`(eq "1" "2")`, `0`},
		{`(ne 1 1)`, `0`},
		{`(ne 1 2)`, `1`},
		{`(ne "1" "1")`, `0`},
		{`(ne "1" "2")`, `1`},
		{`(and 1 1)`, `1`},
		{`(and 1 0)`, `0`},
		{`(and 2 3)`, `1`},
		{`(or 0 1)`, `1`},
		{`(or 0 0)`, `0`},
		{`(or 2 3)`, `1`},
		{`(not 1)`, `0`},
		{`(not 2)`, `0`},
		{`(not 0)`, `1`},
		{`(concat "hello" "world")`, `helloworld`},
		{`(id (void))`, `3`},
		{`(eq (id (void)) (id (void)))`, `1`},
		{`(eq (id 1) (id 2))`, `0`},
		{`(eq (id "str") (id "str"))`, `0`},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			state := NewState(lexAndParse(t, test.input), conf.New())
			assert.Nil(t, state.Execute())
			assert.Equal(t, test.value, state.Value().String())
		})
	}
}
