package runtime

import (
	"fmt"
	"reflect"
	"strconv"
	"sync/atomic"

	"github.com/gogim1/goscript/ast"
	"github.com/gogim1/goscript/file"
)

type Value interface {
	GetId() int64
	SetId(int64)
	String() string
}

type Base struct {
	Location int
	Id       int64
}

func (n *Base) GetId() int64 {
	return n.Id
}

func (n *Base) SetId(i int64) {
	n.Id = i
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

func (v *Number) lt(other *Number) bool {
	lhs := v.Numerator * other.Denominator
	rhs := v.Denominator * other.Numerator
	return lhs < rhs
}

var (
	trueValue  = NewNumber(1, 1)
	falseValue = NewNumber(0, 1)
)

type String struct {
	Base
	Value string
}

func (v *String) String() string {
	return v.Value
}

type Closure struct {
	Base
	Env []envItem
	Fun *ast.LambdaNode
}

func (v *Closure) String() string {
	return fmt.Sprintf("<closure evaluated at %s>", v.Fun.Location)
}

type Continuation struct {
	Base
	SourceLocation file.SourceLocation
	Stack          []*layer
}

func (v *Continuation) String() string {
	return fmt.Sprintf("<continuation evaluated at %s>", v.SourceLocation)
}

var globalId int64 = 0

func NewVoid() *Void {
	ret := &Void{}
	id := atomic.AddInt64(&globalId, 1)
	ret.SetId(id)
	return ret
}

func NewNumber(n, d int) *Number {
	ret := &Number{
		Numerator:   n,
		Denominator: d,
	}
	id := atomic.AddInt64(&globalId, 1)
	ret.SetId(id)
	return ret
}

func NewString(v string) *String {
	ret := &String{
		Value: v,
	}
	id := atomic.AddInt64(&globalId, 1)
	ret.SetId(id)
	return ret
}

func NewClosure(env []envItem, fun *ast.LambdaNode) *Closure {
	ret := &Closure{
		Env: env,
		Fun: fun,
	}
	id := atomic.AddInt64(&globalId, 1)
	ret.SetId(id)
	return ret
}

func NewContinuation(sl file.SourceLocation, stack []*layer) *Continuation {
	ret := &Continuation{
		SourceLocation: sl,
		Stack:          stack,
	}
	id := atomic.AddInt64(&globalId, 1)
	ret.SetId(id)
	return ret
}
