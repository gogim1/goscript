package runtime

import (
	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/parser"
)

type layer struct {
	expr parser.ExprNode
	pc   int
}

type state struct {
	value Value
	stack []layer
}

func NewState(expr parser.ExprNode) *state {
	return &state{
		stack: []layer{
			{expr: nil},
			{expr: expr, pc: 0},
		},
	}
}

func (s *state) Value() string {
	return s.value.String()
}

func (s *state) Execute() *file.Error {
	for {
		layer := s.stack[len(s.stack)-1]
		if layer.expr == nil {
			return nil
		}

		switch expr := layer.expr.(type) {
		case *parser.NumberNode:
			s.value = NewNumberValue(expr.Numerator, expr.Denominator)
			s.stack = s.stack[:len(s.stack)-1]
		default:
			return &file.Error{Location: expr.GetLocation(), Message: "unrecognized AST node"}
		}
	}
}
