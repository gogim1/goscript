package parser

import (
	"math"
	"unicode"

	"github.com/gogim1/goscript/file"
)

type ExprNode interface {
	GetLocation() file.SourceLocation
	SetLocation(file.SourceLocation)
}

type base struct {
	location file.SourceLocation
}

func (n *base) GetLocation() file.SourceLocation {
	return n.location
}

func (n *base) SetLocation(sl file.SourceLocation) {
	n.location = sl
}

type NumberNode struct {
	base
	numerator   int
	denominator int
}

func NewNumberNode(sl file.SourceLocation, n, d int) *NumberNode {
	g := gcd(int(math.Abs(float64(n))), d)
	node := &NumberNode{
		base:        base{location: sl},
		numerator:   n / g,
		denominator: d / g,
	}
	return node
}

type StringNode struct {
	base
	value string
}

func NewStringNode(sl file.SourceLocation, v string) *StringNode {
	node := &StringNode{
		base:  base{location: sl},
		value: v,
	}
	return node
}

type IntrinsicNode struct {
	base
	name string
}

func NewIntrinsicNode(sl file.SourceLocation, n string) *IntrinsicNode {
	node := &IntrinsicNode{
		base: base{location: sl},
		name: n,
	}
	return node
}

type VariableNode struct {
	base
	name string
}

func (n *VariableNode) IsLex() bool {
	return unicode.IsLower([]rune(n.name)[0])
}

func (n *VariableNode) IsDyn() bool {
	return unicode.IsUpper([]rune(n.name)[0])
}

func NewVariableNode(sl file.SourceLocation, n string) *VariableNode {
	node := &VariableNode{
		base: base{location: sl},
		name: n,
	}
	return node
}

type LambdaNode struct {
	base
	varList []ExprNode
	expr    ExprNode
}

func NewLambdaNode(sl file.SourceLocation, vl []ExprNode, e ExprNode) *LambdaNode {
	node := &LambdaNode{
		base:    base{location: sl},
		varList: vl,
		expr:    e,
	}
	return node
}

type LetrecVarExprItem struct {
	variable VariableNode
	expr     ExprNode
}

func NewLetrecVarExprItem(v VariableNode, e ExprNode) *LetrecVarExprItem {
	node := &LetrecVarExprItem{
		variable: v,
		expr:     e,
	}
	return node
}

type LetrecNode struct {
	base
	varExprList []LetrecVarExprItem
	expr        ExprNode
}

func NewLetrecNode(sl file.SourceLocation, vel []LetrecVarExprItem, e ExprNode) *LetrecNode {
	node := &LetrecNode{
		base:        base{location: sl},
		varExprList: vel,
		expr:        e,
	}
	return node
}

type IfNode struct {
	base
	cond    ExprNode
	branch1 ExprNode
	branch2 ExprNode
}

func NewIfNode(sl file.SourceLocation, c, b1, b2 ExprNode) *IfNode {
	node := &IfNode{
		base:    base{location: sl},
		cond:    c,
		branch1: b1,
		branch2: b2,
	}
	return node
}

type CallNode struct {
	base
	callee  ExprNode
	argList []ExprNode
}

func NewCallNode(sl file.SourceLocation, c ExprNode, al []ExprNode) *CallNode {
	node := &CallNode{
		base:    base{location: sl},
		callee:  c,
		argList: al,
	}
	return node
}

type SequenceNode struct {
	base
	exprList []ExprNode
}

func NewSequenceNode(sl file.SourceLocation, el []ExprNode) *SequenceNode {
	node := &SequenceNode{
		base:     base{location: sl},
		exprList: el,
	}
	return node
}

type QueryNode struct {
	base
	variable VariableNode
	exprBox  []ExprNode
}

func NewQueryNode(sl file.SourceLocation, v VariableNode, eb []ExprNode) *QueryNode {
	node := &QueryNode{
		base:     base{location: sl},
		variable: v,
		exprBox:  eb,
	}
	return node
}

type AccessNode struct {
	base
	variable VariableNode
	expr     ExprNode
}

func NewAccessNode(sl file.SourceLocation, v VariableNode, e ExprNode) *AccessNode {
	node := &AccessNode{
		base:     base{location: sl},
		variable: v,
		expr:     e,
	}
	return node
}
