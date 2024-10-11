package runtime

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gogim1/goscript/ast"
	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
	"github.com/gogim1/goscript/parser"
)

func (s *state) preserve() []*layer {
	layers := make([]*layer, len(s.stack))
	deepcopy(&layers, s.stack)
	return layers
}

func (s *state) restore(layers []*layer) {
	deepcopy(&s.stack, layers)
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

func (s *state) VisitNumberNode(n *ast.NumberNode) *file.Error {
	s.value = NewNumber(n.Numerator, n.Denominator)
	s.stack = s.stack[:len(s.stack)-1]
	return nil
}

func (s *state) VisitStringNode(n *ast.StringNode) *file.Error {
	s.value = NewString(n.Value)
	s.stack = s.stack[:len(s.stack)-1]
	return nil
}

func (s *state) VisitIntrinsicNode(n *ast.IntrinsicNode) *file.Error {
	l := s.stack[len(s.stack)-1]

	switch n.Name {
	case "void":
		s.value = NewVoid()
	case "isVoid":
		if _, ok := l.args[0].(*Void); ok {
			s.value = NewNumber(1, 1)
		} else {
			s.value = NewNumber(0, 1)
		}
	case "isNum":
		if _, ok := l.args[0].(*Number); ok {
			s.value = NewNumber(1, 1)
		} else {
			s.value = NewNumber(0, 1)
		}
	case "isStr":
		if _, ok := l.args[0].(*String); ok {
			s.value = NewNumber(1, 1)
		} else {
			s.value = NewNumber(0, 1)
		}
	case "isCont":
		if _, ok := l.args[0].(*Continuation); ok {
			s.value = NewNumber(1, 1)
		} else {
			s.value = NewNumber(0, 1)
		}
	case "put":
		output := ""
		for _, v := range l.args {
			output += fmt.Sprint(v)
		}
		fmt.Print(output)
		s.value = NewVoid()
	case "getline":
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			s.value = NewVoid()
		} else {
			s.value = NewString(scanner.Text())
		}
	case "eval":
		v, err := run(l.args[0].(*String).String())
		if err != nil {
			return err
		} else {
			s.value = v
		}
	case "callCC":
		s.stack = s.stack[:len(s.stack)-1]

		contStack := s.preserve()
		addr := s.new(NewContinuation(n.GetLocation(), contStack))
		closure := l.args[0].(*Closure)

		env := make([]envItem, len(closure.Env))
		copy(env, closure.Env)
		env = append(env, envItem{closure.Fun.VarList[0].Name, addr})

		s.stack = append(s.stack, &layer{
			env:   &env,
			frame: true,
			expr:  closure.Fun.Expr,
		})
		return nil
	case "reg":
		addr := l.args[1].GetLocation()
		if addr == -1 {
			addr = s.new(l.args[1])
		}
		*(s.stack[0].env) = append(*(s.stack[0].env), envItem{
			name:     l.args[0].(*String).Value,
			location: addr,
		})
		s.value = NewVoid()
	case "go":
		name := l.args[0].(*String).Value
		args := l.args[1:]
		if f, ok := s.ffi[name]; !ok {
			return &file.Error{
				Location: n.GetLocation(),
				Message:  "FFI encountered unregistered function",
			}
		} else {
			s.value = f(args...)
		}
	default:
		return &file.Error{
			Location: n.GetLocation(),
			Message:  "unrecognized intrinsic function call",
		}
	}
	s.stack = s.stack[:len(s.stack)-1]
	return nil
}

func (s *state) VisitVariableNode(n *ast.VariableNode) *file.Error {
	l := s.stack[len(s.stack)-1]
	location := -1
	if n.Kind == ast.Lexical {
		location = lookupEnv(n.Name, *l.env)
	} else {
		location = lookupStack(n.Name, s.stack)
	}
	if location == -1 {
		return &file.Error{
			Location: n.GetLocation(),
			Message:  "undefined variable",
		}
	}
	s.value = s.heap[location]
	s.stack = s.stack[:len(s.stack)-1]
	return nil
}

func (s *state) VisitLambdaNode(n *ast.LambdaNode) *file.Error {
	l := s.stack[len(s.stack)-1]
	s.value = NewClosure(filterLexical(*l.env), n)
	s.stack = s.stack[:len(s.stack)-1]
	return nil
}

func (s *state) VisitLetrecNode(n *ast.LetrecNode) *file.Error {
	l := s.stack[len(s.stack)-1]
	if 1 < l.pc && l.pc <= len(n.VarExprList)+1 {
		v := n.VarExprList[l.pc-2].Variable
		lastLocation := lookupEnv(v.Name, *l.env)
		if lastLocation == -1 {
			panic("this should not happened. panic for testing.") // TODO
		}
		s.heap[lastLocation] = s.value
	}
	if l.pc == 0 {
		for _, ve := range n.VarExprList {
			*l.env = append(*l.env, envItem{
				name:     ve.Variable.Name,
				location: s.new(NewVoid()),
			})
		}
		l.pc++
	} else if l.pc <= len(n.VarExprList) {
		s.stack = append(s.stack, &layer{
			env:  l.env,
			expr: n.VarExprList[l.pc-1].Expr,
		})
		l.pc++
	} else if l.pc == len(n.VarExprList)+1 {
		s.stack = append(s.stack, &layer{
			env:  l.env,
			expr: n.Expr,
		})
		l.pc++
	} else {
		*l.env = (*l.env)[:len(*l.env)-len(n.VarExprList)]
		s.stack = s.stack[:len(s.stack)-1]
	}
	return nil
}

func (s *state) VisitIfNode(n *ast.IfNode) *file.Error {
	l := s.stack[len(s.stack)-1]
	if l.pc == 0 {
		s.stack = append(s.stack, &layer{
			env:  l.env,
			expr: n.Cond,
		})
		l.pc++
	} else if l.pc == 1 {
		if v, ok := s.value.(*Number); !ok {
			return &file.Error{
				Location: n.Cond.GetLocation(),
				Message:  "wrong condition type",
			}
		} else {
			newLayer := &layer{
				env: l.env,
			}
			if v.Numerator != 0 {
				newLayer.expr = n.Branch1
			} else {
				newLayer.expr = n.Branch2
			}
			s.stack = append(s.stack, newLayer)
		}
		l.pc++
	} else {
		s.stack = s.stack[:len(s.stack)-1]
	}
	return nil
}

func (s *state) VisitCallNode(n *ast.CallNode) *file.Error {
	l := s.stack[len(s.stack)-1]
	if callee, ok := n.Callee.(*ast.IntrinsicNode); ok {
		if 1 < l.pc && l.pc <= len(n.ArgList)+1 {
			l.args = append(l.args, s.value)
		}
		if l.pc == 0 {
			l.pc++
		} else if l.pc <= len(n.ArgList) {
			s.stack = append(s.stack, &layer{
				env:  l.env,
				expr: n.ArgList[l.pc-1],
			})
			l.pc++
		} else {
			return callee.Accept(s)
		}
	} else {
		if 2 < l.pc && l.pc <= len(n.ArgList)+2 {
			l.args = append(l.args, s.value)
		}
		if l.pc == 0 {
			s.stack = append(s.stack, &layer{
				env:  l.env,
				expr: n.Callee,
			})
			l.pc++
		} else if l.pc == 1 {
			l.callee = s.value
			l.pc++
		} else if l.pc <= len(n.ArgList)+1 {
			s.stack = append(s.stack, &layer{
				env:  l.env,
				expr: n.ArgList[l.pc-2],
			})
			l.pc++
		} else if l.pc == len(n.ArgList)+2 {
			if closure, ok := l.callee.(*Closure); ok {
				if len(l.args) != len(closure.Fun.VarList) {
					return &file.Error{
						Location: n.GetLocation(),
						Message:  "wrong number of arguments given to callee",
					}
				}
				env := make([]envItem, len(closure.Env))
				copy(env, closure.Env)

				for i, v := range closure.Fun.VarList {
					addr := l.args[i].GetLocation()
					if addr == -1 {
						addr = s.new(l.args[i])
					}
					env = append(env, envItem{
						name:     v.Name,
						location: addr,
					})
				}
				s.stack = append(s.stack, &layer{
					env:   &env,
					frame: true,
					expr:  closure.Fun.Expr,
				})
				l.pc++
			} else if continuation, ok := l.callee.(*Continuation); ok {
				s.restore(continuation.Stack)
				// ?
			} else {
				return &file.Error{
					Location: n.Callee.GetLocation(),
					Message:  "calling non-callable object",
				}
			}
		} else {
			s.stack = s.stack[:len(s.stack)-1]
		}
	}
	return nil
}

func (s *state) VisitSequenceNode(n *ast.SequenceNode) *file.Error {
	l := s.stack[len(s.stack)-1]
	if l.pc < len(n.ExprList) {
		s.stack = append(s.stack, &layer{
			env:  l.env,
			expr: n.ExprList[l.pc],
		})
		l.pc++
	} else {
		s.stack = s.stack[:len(s.stack)-1]
	}
	return nil
}

func (s *state) VisitAccessNode(n *ast.AccessNode) *file.Error {
	return nil
}

func (s *state) VisitQueryNode(n *ast.QueryNode) *file.Error {
	return nil
}
