package runtime

import (
	"unicode"

	"github.com/gogim1/goscript/ast"
	"github.com/gogim1/goscript/file"
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
	store []Value
}

func NewState(expr ast.ExprNode) *state {
	return &state{
		stack: []*layer{
			{expr: nil},
			{expr: expr},
		},
	}
}

func (s *state) Value() string {
	return s.value.String()
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
	location := len(s.store)
	s.store = append(s.store, value)
	value.SetLocation(location)
	return location
}