package main

import (
	"fmt"

	"github.com/gogim1/goscript/conf"
	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
	"github.com/gogim1/goscript/parser"
	"github.com/gogim1/goscript/runtime"
)

// call goscript function from golang
func func1() {
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
		fmt.Println(err)
		return
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		fmt.Println(err)
		return
	}
	state := runtime.NewState(node, conf.New())
	err = state.Execute()
	if err != nil {
		fmt.Println(err)
		return
	}
	v, err := state.Call("test0")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)

	_, err = state.Call("test1", 42)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func concat(args ...runtime.Value) runtime.Value {
	s := ""
	for _, v := range args {
		s += v.(*runtime.String).Value
	}
	return runtime.NewString(s)
}

// call golang function from goscript
func func2() {
	tokens, err := lexer.Lex(file.NewSource(`
letrec (s = (go "concat" "hello" " " "world")) {
  (put s "\n")
}
	`))
	if err != nil {
		fmt.Println(err)
		return
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = runtime.NewState(node, conf.New()).Register("concat", concat).Execute()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	func1()
	func2()
}
