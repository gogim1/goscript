package main

import (
	"fmt"

	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
	"github.com/gogim1/goscript/parser"
	"github.com/gogim1/goscript/runtime"
)

func main() {
	tokens, err := lexer.Lex(file.NewSource(`
	letrec (
		v = 1
	) {
		[
			(reg "test0" lambda(){ v })
			(reg "test1" lambda(v){ (put v "\n") })
		]
	}
	
	`))
	if err != nil {
		fmt.Print(err)
		return
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		fmt.Print(err)
		return
	}
	state := runtime.NewState(node)
	err = state.Execute()
	if err != nil {
		fmt.Print(err)
		return
	}
	v, err := state.CallScriptFunction("test0", nil)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println(v)

	_, err = state.CallScriptFunction("test1", []any{42})
	if err != nil {
		fmt.Print(err)
		return
	}
}
