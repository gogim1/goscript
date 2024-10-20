package main

import (
	"fmt"

	"github.com/gogim1/goscript/conf"
	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
	"github.com/gogim1/goscript/parser"
	"github.com/gogim1/goscript/runtime"
)

func main() {
	tokens, err := lexer.Lex(file.NewSource(`
	letrec (
		loop = lambda (s) {
		  letrec (
			line = [
			  (put "> ") 
			  (getline)
			  ]
		  ) {
			if (ne line "exit") then [
			  (put (eval line) "\n")
			  (loop (concat (concat (concat s "> ") line) "\n"))
			  ]
			else (put "Bye. Your input history:\n" s)
		  }
		}
	  ) {
		(loop "")
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
	err = runtime.NewState(node, conf.New()).Execute()
	if err != nil {
		fmt.Println(err)
		return
	}
}
