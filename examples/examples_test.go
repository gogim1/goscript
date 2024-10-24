package examples_test

import (
	"fmt"
	"os"

	"github.com/gogim1/goscript/conf"
	"github.com/gogim1/goscript/examples"
	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
	"github.com/gogim1/goscript/parser"
	"github.com/gogim1/goscript/runtime"
)

func run(filepath string, conf *conf.Config) error {
	bytes, e := os.ReadFile(filepath)
	if e != nil {
		return e
	}

	tokens, err := lexer.Lex(file.NewSource(string(bytes)))
	if err != nil {
		return err
	}

	node, err := parser.Parse(tokens)
	if err != nil {
		return err
	}

	state := runtime.NewState(node, conf)
	if err = state.Execute(); err != nil {
		return err
	}
	return nil
}

var gConf = conf.New(
	conf.SetGCTrigger(func() bool { return true }),
	conf.EnableTCO(true),
)

func Example_binaryTree() {
	err := run("./binary-tree.gs", gConf)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func Example_coroutines() {
	err := run("./coroutines.gs", gConf)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	// Output:
	// main
	// task 1
	// main
	// task 2
	// main
	// task 3
}

func Example_exception() {
	err := run("./exception.gs", gConf)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	// Output:
	// enter `try` block
	// enter `catch` block
	// Exception: message
	// exit `catch` block
	// enter `finally` block
	// exit `finally` block
	// enter `try` block
	// exit `try` block
	// enter `finally` block
	// exit `finally` block
}

func Example_list() {
	err := run("./list.gs", gConf)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	// Output:
	// 5
	// 4
	// 3
	// 2
	// 1
}

func Example_multi_stage() {
	err := run("./multi-stage.gs", gConf)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	// Output:
	// EVAL
	// hello world
	// hello world
}

func Example_mutual_recursion() {
	err := run("./mutual-recursion.gs", gConf)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	// Output:
	// 5
	// 4
	// 3
	// 2
	// 1
	// 0
}

func Example_oop() {
	err := run("./oop.gs", gConf)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	// Output:
	// 1
	// 2
	// 100
	// 2
}

func Example_scope() {
	err := run("./scope.gs", gConf)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	// Output:
	// 1
	// 303
}

func Example_values() {
	err := run("./values.gs", gConf)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	// Output:
	// <void>
	// 1/2
	// str
	// <closure evaluated at (SourceLocation 5 7)>
	// <continuation evaluated at (SourceLocation 6 7)>
}

func Example_y_combinator() {
	err := run("./y-combinator.gs", gConf)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	// Output:
	// 1 120 3628800
}

func Example_call_goscript() {
	examples.CallGoscript()

	// Output:
	// 1
	// 42

}

func Example_call_golang() {
	examples.CallGolang()

	// Output:
	// hello world
}
