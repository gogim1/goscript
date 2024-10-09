package runtime

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gogim1/goscript/ast"
	"github.com/gogim1/goscript/file"
)

func (s *state) VisitNumberNode(n *ast.NumberNode) *file.Error {
	s.value = &Number{
		Numerator:   n.Numerator,
		Denominator: n.Denominator,
	}
	s.value.SetLocation(-1)
	s.stack = s.stack[:len(s.stack)-1]
	return nil
}

func (s *state) VisitStringNode(n *ast.StringNode) *file.Error {
	s.value = &String{
		Value: n.Value,
	}
	s.value.SetLocation(-1)
	s.stack = s.stack[:len(s.stack)-1]
	return nil
}

func (s *state) VisitIntrinsicNode(n *ast.IntrinsicNode) *file.Error {
	l := s.stack[len(s.stack)-1]

	switch n.Name {
	case "void":
		s.value = &Void{}
	case "isVoid":
		if _, ok := l.args[0].(*Void); ok {
			s.value = &Number{Numerator: 1, Denominator: 1}
		} else {
			s.value = &Number{Numerator: 0, Denominator: 1}
		}
	case "isNum":
		if _, ok := l.args[0].(*Number); ok {
			s.value = &Number{Numerator: 1, Denominator: 1}
		} else {
			s.value = &Number{Numerator: 0, Denominator: 1}
		}
	case "isStr":
		if _, ok := l.args[0].(*String); ok {
			s.value = &Number{Numerator: 1, Denominator: 1}
		} else {
			s.value = &Number{Numerator: 0, Denominator: 1}
		}
	case "isCont":
		if _, ok := l.args[0].(*Continuation); ok {
			s.value = &Number{Numerator: 1, Denominator: 1}
		} else {
			s.value = &Number{Numerator: 0, Denominator: 1}
		}
	case "put":
		output := ""
		for _, v := range l.args {
			output += fmt.Sprint(v)
		}
		fmt.Print(output)
		s.value = &Void{}
	case "getline":
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			s.value = &Void{}
		} else {
			s.value = &String{Value: scanner.Text()}
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
		contStack := s.cloneStack()
		addr := s.new(&Continuation{
			SourceLocation: n.GetLocation(),
			Stack:          contStack,
		})
		closure := l.args[0].(*Closure)
		s.stack = append(s.stack, &layer{
			env:  append(closure.Env, envItem{closure.Fun.VarList[0].Name, addr}),
			expr: closure.Fun.Expr,
		})
		return nil
	case "reg":
		addr := l.args[1].GetLocation()
		if addr == -1 {
			addr = s.new(l.args[1])
		}
		s.stack[0].env = append(s.stack[0].env, envItem{
			name:     l.args[0].(*String).Value,
			location: addr,
		})
		s.value = &Void{}
	case "go":
		goFunName := l.args[0].(*String).Value
		goArgs := l.args[1:]
		if f, ok := s.goFunctions[goFunName]; !ok {
			return &file.Error{
				Location: n.GetLocation(),
				Message:  "FFI encountered unregistered function",
			}
		} else {
			s.value = f(goArgs)
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
		location = lookupEnv(n.Name, l.env)
	} else {
		location = lookupStack(n.Name, s.stack)
	}
	if location == -1 {
		return &file.Error{
			Location: n.GetLocation(),
			Message:  "undefined variable",
		}
	}
	s.value = s.store[location]
	s.stack = s.stack[:len(s.stack)-1]
	return nil
}

func (s *state) VisitLambdaNode(n *ast.LambdaNode) *file.Error {
	l := s.stack[len(s.stack)-1]
	s.value = &Closure{
		Env: s.filterLexical(l.env),
		Fun: n,
	}
	s.value.SetLocation(-1)
	s.stack = s.stack[:len(s.stack)-1]
	return nil
}

func (s *state) VisitLetrecNode(n *ast.LetrecNode) *file.Error {
	l := s.stack[len(s.stack)-1]
	if 1 < l.pc && l.pc <= len(n.VarExprList)+1 {
		v := n.VarExprList[l.pc-2].Variable
		lastLocation := lookupEnv(v.Name, l.env)
		if lastLocation == -1 {
			panic("this should not happened. panic for testing.") // TODO
		}
		s.store[lastLocation] = s.value
	}
	if l.pc == 0 {
		for _, ve := range n.VarExprList {
			l.env = append(l.env, envItem{
				name:     ve.Variable.Name,
				location: s.new(&Void{}),
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
				for i, v := range closure.Fun.VarList {
					addr := l.args[i].GetLocation()
					if addr == -1 {
						addr = s.new(l.args[i])

					}
					closure.Env = append(closure.Env, envItem{
						name:     v.Name,
						location: addr,
					})
				}
				s.stack = append(s.stack, &layer{
					env:  closure.Env,
					expr: closure.Fun.Expr,
				})
				l.pc++
			} else if continuation, ok := l.callee.(*Continuation); ok {
				s.restoreStack(continuation.Stack)
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
