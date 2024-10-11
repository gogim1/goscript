package runtime

import (
	"unicode"

	"github.com/gogim1/goscript/ast"
	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
	"github.com/gogim1/goscript/parser"
)

type layer struct {
	env    []envItem
	expr   ast.ExprNode
	pc     int
	args   []Value
	callee Value
}

type state struct {
	value Value
	stack []*layer
	heap  []Value
	ffi   map[string]func(...Value) Value
}

func NewState(expr ast.ExprNode) *state {
	return &state{
		stack: []*layer{
			{expr: nil},
			{expr: expr},
		},
		ffi: make(map[string]func(...Value) Value),
	}
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
	}
}

func lookupEnv(name string, env []envItem) int {
	for i := len(env) - 1; i >= 0; i-- {
		if env[i].name == name {
			return env[i].location
		}
	}
	return -1
}

func lookupStack(name string, layers []*layer) int {
	for i := len(layers) - 1; i >= 0; i-- {
		l := layers[i]
		for j := len(l.env) - 1; j >= 0; j-- {
			if l.env[j].name == name {
				return l.env[j].location
			}
		}
	}
	return -1
}

func (s *state) filterLexical(env []envItem) []envItem {
	newEnv := []envItem{}
	for _, item := range env {
		if len(item.name) > 0 && unicode.IsLower([]rune(item.name)[0]) {
			newEnv = append(newEnv, item)
		}
	}
	return newEnv
}

func (s *state) new(value Value) int {
	location := len(s.heap)
	s.heap = append(s.heap, value)
	value.SetLocation(location)
	return location
}

func (s *state) cloneStack() []*layer {
	layers := make([]*layer, len(s.stack))
	for i, l := range s.stack {
		layers[i] = &layer{
			env:    l.env,
			expr:   l.expr,
			pc:     l.pc,
			args:   l.args,
			callee: l.callee,
		}
	}

	return layers
}

func (s *state) restoreStack(layers []*layer) {
	s.stack = make([]*layer, len(layers))
	for i, l := range layers {
		s.stack[i] = &layer{
			env:    l.env,
			expr:   l.expr,
			pc:     l.pc,
			args:   l.args,
			callee: l.callee,
		}
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
	err := s.Execute()
	if err != nil {
		return nil, err
	}
	return s.value, nil
}

func (s *state) Register(name string, fun func(...Value) Value) *state {
	s.ffi[name] = fun
	return s
}

func run(src string) (Value, *file.Error) {
	tokens, err := lexer.Lex(file.NewSource(src))
	if err != nil {
		return nil, err
	}

	node, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}

	state := NewState(node)
	if err = state.Execute(); err != nil {
		return nil, err
	}
	return state.Value(), nil
}
