package runtime

import (
	"fmt"
	"os"

	"github.com/gogim1/goscript/ast"
	"github.com/gogim1/goscript/conf"
	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/stdlib"
)

//go:generate sh -c "go run ./helpers > ./constants_generated.go"

type envItem struct {
	name     string
	location int
}

type layer struct {
	env    *[]envItem
	frame  bool
	tail   bool
	expr   ast.ExprNode
	pc     int
	args   []Value
	callee Value
}

type state struct {
	collector
	config *conf.Config
	value  Value
	stack  []*layer
	heap   []Value
	ffi    map[string]func(...Value) Value
}

func NewState(expr ast.ExprNode, config *conf.Config) *state {
	s := &state{
		config: config,
		stack: []*layer{
			{env: new([]envItem), expr: nil, frame: true},
		},
		ffi: make(map[string]func(...Value) Value),
	}
	s.collector.state = s
	s.collector.values = make(map[int64]struct{})
	s.collector.locations = make(map[int]struct{})
	s.collector.relocation = make(map[int]int)

	if config.UseStd {
		for _, filepath := range stdlib.Paths {
			bytes, err := os.ReadFile(filepath)
			if err != nil {
				fmt.Printf("[WARN] cannot find stdlib `%s`\n", filepath)
				continue
			}
			node, e := lexAndParse(string(bytes))
			if e != nil {
				panic(fmt.Sprintf("[ERROR] load stdlib `%s` failed\n", filepath))
			}
			s.load(node)
			if e := s.Execute(); e != nil {
				panic(fmt.Sprintf("[ERROR] load stdlib `%s` failed\n", filepath))
			}
		}
		s.gc()
	}
	s.load(expr)
	return s
}

func (s *state) load(expr ast.ExprNode) {
	env := make([]envItem, len(*(s.stack[0].env)))
	copy(env, *(s.stack[0].env))
	s.stack = append(s.stack, &layer{env: &env, expr: expr, frame: true})
}

func (s *state) Value() Value {
	return s.value
}

func (s *state) Execute() *file.Error {
	for {
		l := s.stack[len(s.stack)-1]
		if l.expr == nil {
			break
		}

		err := l.expr.Accept(s)
		if err != nil {
			return err
		}
		if s.config.GCTrigger() {
			n := s.gc()
			if s.config.EnableDebug {
				fmt.Printf("[DEBUG] GC collect %d cells\n", n)
			}
		}
	}
	if s.config.EnableDebug {
		printMemUsage()
	}
	return nil
}

func (s *state) Call(name string, args ...any) (Value, *file.Error) {
	sl := file.SourceLocation{Line: -1, Col: -1}
	callee := ast.NewVariableNode(sl, name, ast.Unknown) // TODO: scope
	argList := []ast.ExprNode{}
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			argList = append(argList, ast.NewStringNode(sl, v))
		case int:
			argList = append(argList, ast.NewNumberNode(sl, v, 1))
		default:
			return nil, &file.Error{
				Location: sl,
				Message:  "golang can only use string/int as arguments when calling goscript functions",
			}
		}
	}
	env := make([]envItem, len(*(s.stack[0].env)))
	copy(env, *(s.stack[0].env))

	s.stack = append(s.stack, &layer{
		env:   &env,
		frame: true,
		expr:  ast.NewCallNode(sl, callee, argList),
	})

	if err := s.Execute(); err != nil {
		return nil, err
	}
	return s.value, nil
}

func (s *state) Register(name string, fun func(...Value) Value) *state {
	s.ffi[name] = fun
	return s
}
