package runtime

import (
	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/parser"
)

type layer struct {
	env  []envItem
	expr parser.ExprNode
	pc   int
}

type state struct {
	value Value
	stack []*layer
	store []Value
}

func NewState(expr parser.ExprNode) *state {
	return &state{
		stack: []*layer{
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
		l := s.stack[len(s.stack)-1]
		if l.expr == nil {
			return nil
		}

		switch expr := l.expr.(type) {
		case *parser.NumberNode:
			s.value = NewNumberValue(expr.Numerator, expr.Denominator)
			s.stack = s.stack[:len(s.stack)-1]
		case *parser.StringNode:
			s.value = NewStringValue(expr.Value.String())
			s.stack = s.stack[:len(s.stack)-1]
		case *parser.LetrecNode:
			if 1 < l.pc && l.pc <= len(expr.VarExprList) {
				v := expr.VarExprList[l.pc-2].Variable
				lastLocation := s.lookupEnv(v.Name.String(), l.env)
				if lastLocation == -1 {
					panic("this should not happened. panic for testing.")
				}
				s.store[lastLocation] = s.value
			}
			if l.pc == 0 {
				for _, ve := range expr.VarExprList {
					l.env = append(l.env, envItem{
						name:     ve.Variable.Name.String(),
						location: s.new(NewVoidValue()),
					})
				}
				l.pc++
			} else if l.pc <= len(expr.VarExprList) {
				s.stack = append(s.stack, &layer{
					env:  l.env,
					expr: expr.VarExprList[l.pc-1].Expr,
				})
				l.pc++
			} else if l.pc == len(expr.VarExprList)+1 {
				s.stack = append(s.stack, &layer{
					env:  l.env,
					expr: expr.Expr,
				})
				l.pc++
			} else {
				l.env = l.env[:len(l.env)-len(expr.VarExprList)]
				s.stack = s.stack[:len(s.stack)-1]
			}
		case *parser.IfNode:
			if l.pc == 0 {
				s.stack = append(s.stack, &layer{
					env:  l.env,
					expr: expr.Cond,
				})
				l.pc++
			} else if l.pc == 1 {
				if n, ok := s.value.(*Number); !ok {
					return &file.Error{
						Location: expr.Cond.GetLocation(),
						Message:  "wrong condition type",
					}
				} else {
					newLayer := &layer{
						env: l.env,
					}
					if n.Numerator != 0 {
						newLayer.expr = expr.Branch1
					} else {
						newLayer.expr = expr.Branch2
					}
					s.stack = append(s.stack, newLayer)
				}
				l.pc++
			} else {
				s.stack = s.stack[:len(s.stack)-1]
			}
		default:
			return &file.Error{Location: expr.GetLocation(), Message: "unrecognized AST node"}
		}
	}
}

func (s *state) lookupEnv(name string, env []envItem) int {
	for i := len(env) - 1; i >= 0; i-- {
		if env[i].name == name {
			return env[i].location
		}
	}
	return -1
}

func (s *state) new(value Value) int {
	location := len(s.store)
	s.store = append(s.store, value)
	value.SetLocation(location)
	return location
}
