package runtime

import (
	"github.com/gogim1/goscript/ast"
	"github.com/gogim1/goscript/conf"
	"github.com/gogim1/goscript/file"
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
			{env: new([]envItem), expr: expr, frame: true},
		},
		ffi: make(map[string]func(...Value) Value),
	}
	s.collector.state = s
	s.collector.values = make(map[int64]struct{})
	s.collector.locations = make(map[int]struct{})
	s.collector.relocation = make(map[int]int)
	return s
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
			s.gc()
		}
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
