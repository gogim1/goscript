package runtime

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gogim1/goscript/ast"
	"github.com/gogim1/goscript/file"
)

type Value interface {
	GetLocation() int
	SetLocation(int)
	String() string
}

type Base struct {
	Location int
}

func (n *Base) GetLocation() int {
	return n.Location
}

func (n *Base) SetLocation(l int) {
	n.Location = l
}

var (
	ValueType        = reflect.TypeOf(new(Value)).Elem()
	VoidType         = reflect.TypeOf(Void{})
	StringType       = reflect.TypeOf(String{})
	NumberType       = reflect.TypeOf(Number{})
	ClosureType      = reflect.TypeOf(Closure{})
	ContinuationType = reflect.TypeOf(Continuation{})
)

type Void struct {
	Base
}

func (v *Void) String() string {
	return "<void>"
}

func NewVoid() *Void {
	ret := &Void{}
	ret.SetLocation(-1)
	return ret
}

type Number struct {
	Base
	Numerator   int
	Denominator int
}

func (v *Number) String() string {
	if v.Denominator == 1 {
		return strconv.Itoa(v.Numerator)
	}
	return strconv.Itoa(v.Numerator) + "/" + strconv.Itoa(v.Denominator)
}

func NewNumber(n, d int) *Number {
	ret := &Number{
		Numerator:   n,
		Denominator: d,
	}
	ret.SetLocation(-1)
	return ret
}

type String struct {
	Base
	Value string
}

func (v *String) String() string {
	return v.Value
}

func NewString(v string) *String {
	ret := &String{
		Value: v,
	}
	ret.SetLocation(-1)
	return ret
}

type Closure struct {
	Base
	Env []envItem
	Fun *ast.LambdaNode
}

func (v *Closure) String() string {
	return fmt.Sprintf("<closure evaluated at %s>", v.Fun.Location)
}

func NewClosure(env []envItem, fun *ast.LambdaNode) *Closure {
	ret := &Closure{
		Env: env,
		Fun: fun,
	}
	ret.SetLocation(-1)
	return ret
}

type Continuation struct {
	Base
	SourceLocation file.SourceLocation
	Stack          []*layer
}

func (v *Continuation) String() string {
	return fmt.Sprintf("<continuation evaluated at %s>", v.SourceLocation)
}

func NewContinuation(sl file.SourceLocation, stack []*layer) *Continuation {
	ret := &Continuation{
		SourceLocation: sl,
		Stack:          stack,
	}
	ret.SetLocation(-1)
	return ret
}
