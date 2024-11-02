package main

import (
	"fmt"
	"os"

	"github.com/gogim1/goscript/conf"
	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
	"github.com/gogim1/goscript/parser"
	"github.com/gogim1/goscript/runtime"
)

var counter = 0

func trigger() bool {
	counter++
	return counter%1000 == 0
}

func main() {
	bytes, err := os.ReadFile("./examples/repl.gs")
	if err != nil {
		fmt.Println(err)
		return
	}
	tokens, _ := lexer.Lex(file.NewSource(string(bytes)))
	node, _ := parser.Parse(tokens)
	e := runtime.NewState(node, conf.New(
		conf.SetGCTrigger(trigger),
		conf.EnableTCO(true),
		conf.EnableDebug(false),
		conf.UseStd(true),
	)).Execute()
	if e != nil {
		fmt.Println(e)
		return
	}
}
