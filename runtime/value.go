package runtime

import (
	"fmt"
	"math"
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

func (v *Number) add(other *Number) *Number {
	n1 := v.Numerator*other.Denominator + other.Numerator*v.Denominator
	d1 := v.Denominator * other.Denominator
	g1 := gcd(int(math.Abs(float64(n1))), d1)
	return NewNumber(n1/g1, d1/g1)
}

func (v *Number) sub(other *Number) *Number {
	n1 := v.Numerator*other.Denominator - other.Numerator*v.Denominator
	d1 := v.Denominator * other.Denominator
	g1 := gcd(int(math.Abs(float64(n1))), d1)
	return NewNumber(n1/g1, d1/g1)
}

func (v *Number) mul(other *Number) *Number {
	n1 := v.Numerator * other.Numerator
	d1 := v.Denominator * other.Denominator
	g1 := gcd(int(math.Abs(float64(n1))), d1)
	return NewNumber(n1/g1, d1/g1)
}

func (v *Number) div(other *Number) *Number {
	n1 := v.Numerator * other.Denominator
	d1 := v.Denominator * other.Numerator
	if d1 < 0 {
		n1 = -n1
		d1 = -d1
	}
	g1 := gcd(int(math.Abs(float64(n1))), d1)
	return NewNumber(n1/g1, d1/g1)
}

func (v *Number) lt(other *Number) bool {
	lhs := v.Numerator * other.Denominator
	rhs := v.Denominator * other.Numerator
	return lhs < rhs
}

func NewNumber(n, d int) *Number {
	ret := &Number{
		Numerator:   n,
		Denominator: d,
	}
	ret.SetLocation(-1)
	return ret
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
