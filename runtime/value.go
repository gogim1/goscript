package runtime

import (
	"fmt"
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

type String struct {
	Base
	Value string
}

func (v *String) String() string {
	return v.Value
}

type envItem struct {
	name     string
	location int
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
	Stack          []layer
}

func (v *Continuation) String() string {
	return fmt.Sprintf("<continuation evaluated at %s>", v.SourceLocation)
}
