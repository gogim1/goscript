package runtime

import (
	"fmt"

	"github.com/gogim1/goscript/ast"
	"github.com/gogim1/goscript/file"
)

type envItem struct {
	name     string
	location int
}

type layer struct {
	env    *[]envItem
	frame  bool
	expr   ast.ExprNode
	pc     int
	args   []Value
	callee Value
}

type state struct {
	collector
	value Value
	stack []*layer
	heap  []Value
	ffi   map[string]func(...Value) Value
}

func NewState(expr ast.ExprNode) *state {
	env := new([]envItem)
	s := &state{
		stack: []*layer{
			{env: env, expr: nil, frame: true},
			{env: env, expr: expr},
		},
		ffi: make(map[string]func(...Value) Value),
	}
	s.collector.state = s
	return s
}

func (s *state) Value() Value {
	return s.value
}

func (s *state) Execute() *file.Error {
	for {
		l := s.stack[len(s.stack)-1]
		if l.expr == nil {
			return nil
		}

		err := l.expr.Accept(s)
		if err != nil {
			return err
		}
		fmt.Printf("[Note] GC collected %d cells\n", s.gc())
	}
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
	s.stack = append(s.stack, &layer{
		env:  s.stack[0].env,
		expr: ast.NewCallNode(sl, callee, argList),
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
