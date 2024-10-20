package runtime

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"

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

func retrieveNumberValue(numerator, denominator int) Value {
	if denominator == 1 && numValue[numerator] != nil {
		return numValue[numerator]
	} else {
		return NewNumber(numerator, denominator)
	}
}

func (s *state) VisitNumberNode(n *ast.NumberNode) *file.Error {
	s.value = retrieveNumberValue(n.Numerator, n.Denominator)
	s.stack = s.stack[:len(s.stack)-1]
	return nil
}

func retrieveStringValue(v string) Value {
	if len(v) == 0 {
		return emptyStrValue
	} else {
		return NewString(v)
	}
}

func (s *state) VisitStringNode(n *ast.StringNode) *file.Error {
	s.value = retrieveStringValue(n.Value)
	s.stack = s.stack[:len(s.stack)-1]
	return nil
}

func (s *state) VisitIntrinsicNode(n *ast.IntrinsicNode) *file.Error {
	l := s.stack[len(s.stack)-1]

	switch n.Name {
	case "void":
		if err := typeCheck(l.expr.GetLocation(), l.args, nil); err != nil {
			s.value = voidValue
			return err
		}
		s.value = voidValue
	case "isvoid":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{ValueType}); err != nil {
			s.value = voidValue
			return err
		}
		if _, ok := l.args[0].(*Void); ok {
			s.value = trueValue
		} else {
			s.value = falseValue
		}
	case "isnum":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{ValueType}); err != nil {
			s.value = voidValue
			return err
		}
		if _, ok := l.args[0].(*Number); ok {
			s.value = trueValue
		} else {
			s.value = falseValue
		}
	case "isstr":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{ValueType}); err != nil {
			s.value = voidValue
			return err
		}
		if _, ok := l.args[0].(*String); ok {
			s.value = trueValue
		} else {
			s.value = falseValue
		}
	case "isclo":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{ValueType}); err != nil {
			s.value = voidValue
			return err
		}
		if _, ok := l.args[0].(*Closure); ok {
			s.value = trueValue
		} else {
			s.value = falseValue
		}
	case "iscont":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{ValueType}); err != nil {
			s.value = voidValue
			return err
		}
		if _, ok := l.args[0].(*Continuation); ok {
			s.value = trueValue
		} else {
			s.value = falseValue
		}
	case "add":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{NumberType, NumberType}); err != nil {
			s.value = voidValue
			return err
		}
		lhs, rhs := l.args[0].(*Number), l.args[1].(*Number)
		n1 := lhs.Numerator*rhs.Denominator + rhs.Numerator*lhs.Denominator
		d1 := lhs.Denominator * rhs.Denominator
		g1 := gcd(int(math.Abs(float64(n1))), d1)
		s.value = retrieveNumberValue(n1/g1, d1/g1)
	case "sub":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{NumberType, NumberType}); err != nil {
			s.value = voidValue
			return err
		}
		lhs, rhs := l.args[0].(*Number), l.args[1].(*Number)
		n1 := lhs.Numerator*rhs.Denominator - rhs.Numerator*lhs.Denominator
		d1 := lhs.Denominator * rhs.Denominator
		g1 := gcd(int(math.Abs(float64(n1))), d1)
		s.value = retrieveNumberValue(n1/g1, d1/g1)
	case "mul":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{NumberType, NumberType}); err != nil {
			s.value = voidValue
			return err
		}
		lhs, rhs := l.args[0].(*Number), l.args[1].(*Number)
		n1 := lhs.Numerator * rhs.Numerator
		d1 := lhs.Denominator * rhs.Denominator
		g1 := gcd(int(math.Abs(float64(n1))), d1)
		s.value = retrieveNumberValue(n1/g1, d1/g1)
	case "div":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{NumberType, NumberType}); err != nil {
			s.value = voidValue
			return err
		}
		lhs, rhs := l.args[0].(*Number), l.args[1].(*Number)
		if rhs.Numerator == 0 {
			s.value = voidValue
			return &file.Error{
				Location: l.expr.GetLocation(),
				Message:  "division by zero",
			}
		}
		n1 := lhs.Numerator * rhs.Denominator
		d1 := lhs.Denominator * rhs.Numerator
		if d1 < 0 {
			n1 = -n1
			d1 = -d1
		}
		g1 := gcd(int(math.Abs(float64(n1))), d1)
		s.value = retrieveNumberValue(n1/g1, d1/g1)
	case "lt":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{NumberType, NumberType}); err != nil {
			s.value = voidValue
			return err
		}
		lhs, rhs := l.args[0].(*Number), l.args[1].(*Number)
		if lhs.lt(rhs) {
			s.value = trueValue
		} else {
			s.value = falseValue
		}
	case "gt":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{NumberType, NumberType}); err != nil {
			s.value = voidValue
			return err
		}
		lhs, rhs := l.args[0].(*Number), l.args[1].(*Number)
		if rhs.lt(lhs) {
			s.value = trueValue
		} else {
			s.value = falseValue
		}
	case "ge":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{NumberType, NumberType}); err != nil {
			s.value = voidValue
			return err
		}
		lhs, rhs := l.args[0].(*Number), l.args[1].(*Number)
		if !lhs.lt(rhs) {
			s.value = trueValue
		} else {
			s.value = falseValue
		}
	case "le":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{NumberType, NumberType}); err != nil {
			s.value = voidValue
			return err
		}
		lhs, rhs := l.args[0].(*Number), l.args[1].(*Number)
		if !rhs.lt(lhs) {
			s.value = trueValue
		} else {
			s.value = falseValue
		}
	case "eq":
		if len(l.args) != 2 || reflect.TypeOf(l.args[0]) != reflect.TypeOf(l.args[1]) ||
			(reflect.TypeOf(l.args[0]).Elem() != NumberType && reflect.TypeOf(l.args[0]).Elem() != StringType) {
			s.value = voidValue
			return &file.Error{
				Location: l.expr.GetLocation(),
				Message:  "wrong number/type of arguments given to eq",
			}
		}
		if lhs, ok := l.args[0].(*Number); ok {
			rhs := l.args[1].(*Number)
			if !lhs.lt(rhs) && !rhs.lt(lhs) {
				s.value = trueValue
			} else {
				s.value = falseValue
			}
		} else {
			lhs, rhs := l.args[0].(*String), l.args[1].(*String)
			if lhs.Value == rhs.Value {
				s.value = trueValue
			} else {
				s.value = falseValue
			}
		}
	case "ne":
		if len(l.args) != 2 || reflect.TypeOf(l.args[0]) != reflect.TypeOf(l.args[1]) ||
			(reflect.TypeOf(l.args[0]).Elem() != NumberType && reflect.TypeOf(l.args[0]).Elem() != StringType) {
			s.value = voidValue
			return &file.Error{
				Location: l.expr.GetLocation(),
				Message:  "wrong number/type of arguments given to eq",
			}
		}
		if lhs, ok := l.args[0].(*Number); ok {
			rhs := l.args[1].(*Number)
			if lhs.lt(rhs) || rhs.lt(lhs) {
				s.value = trueValue
			} else {
				s.value = falseValue
			}
		} else {
			lhs, rhs := l.args[0].(*String), l.args[1].(*String)
			if lhs.Value != rhs.Value {
				s.value = trueValue
			} else {
				s.value = falseValue
			}
		}
	case "and":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{NumberType, NumberType}); err != nil {
			s.value = voidValue
			return err
		}
		if l.args[0].(*Number).Numerator != 0 && l.args[1].(*Number).Numerator != 0 {
			s.value = trueValue
		} else {
			s.value = falseValue
		}
	case "or":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{NumberType, NumberType}); err != nil {
			s.value = voidValue
			return err
		}
		if l.args[0].(*Number).Numerator != 0 || l.args[1].(*Number).Numerator != 0 {
			s.value = trueValue
		} else {
			s.value = falseValue
		}
	case "not":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{NumberType}); err != nil {
			s.value = voidValue
			return err
		}
		if l.args[0].(*Number).Numerator == 0 {
			s.value = trueValue
		} else {
			s.value = falseValue
		}
	case "put":
		if len(l.args) == 0 {
			s.value = voidValue
			return &file.Error{
				Location: l.expr.GetLocation(),
				Message:  "wrong number of arguments given to put",
			}
		}
		output := ""
		for _, v := range l.args {
			output += fmt.Sprint(v)
		}
		fmt.Print(output)
		s.value = voidValue
	case "getline":
		if err := typeCheck(l.expr.GetLocation(), l.args, nil); err != nil {
			s.value = voidValue
			return err
		}
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			s.value = voidValue
		} else {
			s.value = retrieveStringValue(scanner.Text())
		}
	case "quote":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{StringType}); err != nil {
			s.value = voidValue
			return err
		}
		str := l.args[0].(*String).Value
		s.value = retrieveStringValue(strconv.Quote(str))
	case "concat":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{StringType, StringType}); err != nil {
			s.value = voidValue
			return err
		}
		s.value = retrieveStringValue(l.args[0].(*String).Value + (l.args[1].(*String)).Value)
	case "eval":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{StringType}); err != nil {
			s.value = voidValue
			return err
		}
		v, err := run(l.args[0].(*String).String())
		if err != nil {
			return err
		} else {
			s.value = v
		}
	case "callcc":
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{ClosureType}); err != nil {
			s.value = voidValue
			return err
		}
		s.stack = s.stack[:len(s.stack)-1]

		contStack := s.preserve()
		addr := s.new(NewContinuation(l.expr.GetLocation(), contStack))
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
		if err := typeCheck(l.expr.GetLocation(), l.args, []reflect.Type{StringType, ClosureType}); err != nil {
			s.value = voidValue
			return err
		}
		*(s.stack[0].env) = append([]envItem{{
			name:     l.args[0].(*String).Value,
			location: s.new(l.args[1]),
		}}, *(s.stack[0].env)...)
		s.value = voidValue
	case "go":
		if len(l.args) == 0 || reflect.TypeOf(l.args[0]).Elem() != StringType {
			s.value = voidValue
			return &file.Error{
				Location: l.expr.GetLocation(),
				Message:  "FFI expects a string (Golang function name) as the first argument",
			}
		}
		name := l.args[0].(*String).Value
		args := l.args[1:]
		if f, ok := s.ffi[name]; !ok {
			s.value = voidValue
			return &file.Error{
				Location: l.expr.GetLocation(),
				Message:  "FFI encountered unregistered function",
			}
		} else {
			s.value = f(args...)
		}
	default:
		s.value = voidValue
		return &file.Error{
			Location: l.expr.GetLocation(),
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
		s.value = voidValue
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
			panic("this should not happened. panic for testing.")
		}
		s.heap[lastLocation] = s.value
	}
	if l.pc == 0 {
		for _, ve := range n.VarExprList {
			*l.env = append(*l.env, envItem{
				name:     ve.Variable.Name,
				location: s.new(voidValue),
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
			s.value = voidValue
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
					s.value = voidValue
					return &file.Error{
						Location: n.GetLocation(),
						Message:  "wrong number of arguments given to callee",
					}
				}
				env := make([]envItem, len(closure.Env))
				copy(env, closure.Env)

				for i, v := range closure.Fun.VarList {
					env = append(env, envItem{
						name:     v.Name,
						location: s.new(l.args[i]),
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
			} else {
				s.value = voidValue
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
	l := s.stack[len(s.stack)-1]
	if l.pc == 0 {
		s.stack = append(s.stack, &layer{
			env:  l.env,
			expr: n.Expr,
		})
		l.pc++
	} else {
		if closure, ok := s.value.(*Closure); ok {
			location := lookupEnv(n.Variable.Name, closure.Env)
			if location == -1 {
				s.value = voidValue
				return &file.Error{
					Location: n.GetLocation(),
					Message:  "undefined variable",
				}
			}
			s.value = s.heap[location]
			s.stack = s.stack[:len(s.stack)-1]
		} else {
			s.value = voidValue
			return &file.Error{
				Location: n.GetLocation(),
				Message:  "lexical variable access applied to non-closure type",
			}
		}
	}
	return nil
}
