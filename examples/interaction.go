package main

import "fmt"

// // call goscript function from golang
// func func1() {
// 	tokens, err := lexer.Lex(file.NewSource(`
// letrec (
// 	v = 1
// ) {
// 	[
// 		(reg "test0" lambda(){ v })
// 		(reg "test1" lambda(v){ (put v "\n") })
// 	]
// }
// 	`))
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	node, err := parser.Parse(tokens)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	state := runtime.NewState(node)
// 	err = state.Execute()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	v, err := state.Call("test0")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(v)

// 	_, err = state.Call("test1", 42)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

// func concat(args ...runtime.Value) runtime.Value {
// 	s := ""
// 	for _, v := range args {
// 		s += v.(*runtime.String).Value
// 	}
// 	return &runtime.String{Value: s}
// }

// // call golang function from goscript
// func func2() {
// 	tokens, err := lexer.Lex(file.NewSource(`
// letrec (s = (go "concat" "hello" " " "world")) {
// 	(put s "\n")
// }
// 	`))
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	node, err := parser.Parse(tokens)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	err = runtime.NewState(node).Register("concat", concat).Execute()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

type envItem struct {
	name     string
	location int
}
type layer struct {
	env   *[]envItem
	frame bool
}

func deepcopy(dst *[]*layer, src []*layer) {
	*dst = make([]*layer, len(src))
	for i, l := range src {
		(*dst)[i] = &layer{
			frame: l.frame,
		}
		env := make([]envItem, len(*l.env))
		copy(env, *l.env)
		(*dst)[i].env = &env
	}
}
func main() {
	// func1()
	// func2()
	env := new([]envItem)
	list := []*layer{{env: env, frame: true}}
	dst := []*layer{}
	deepcopy(&dst, list)
	fmt.Printf("%+v\n", *dst[0])
}
