package ast

import (
	"math"

	"github.com/gogim1/goscript/file"
)

type ExprNode interface {
	GetLocation() file.SourceLocation
	Accept(Visitor) *file.Error
}

type Base struct {
	Location file.SourceLocation
}

func (n *Base) GetLocation() file.SourceLocation {
	return n.Location
}

type NumberNode struct {
	Base
	Numerator   int
	Denominator int
}

func NewNumberNode(sl file.SourceLocation, n, d int) *NumberNode {
	g := gcd(int(math.Abs(float64(n))), d)
	node := &NumberNode{
		Base:        Base{Location: sl},
		Numerator:   n / g,
		Denominator: d / g,
	}
	return node
}

type StringNode struct {
	Base
	Value string
}

func NewStringNode(sl file.SourceLocation, v string) *StringNode {
	node := &StringNode{
		Base:  Base{Location: sl},
		Value: v,
	}
	return node
}

type IntrinsicNode struct {
	Base
	Name string
}

func NewIntrinsicNode(sl file.SourceLocation, n string) *IntrinsicNode {
	node := &IntrinsicNode{
		Base: Base{Location: sl},
		Name: n,
	}
	return node
}

type ScopeKind int

const (
	Unknown ScopeKind = iota
	Lexical
	Dynamic
)

type VariableNode struct {
	Base
	Name string
	Kind ScopeKind
}

func NewVariableNode(sl file.SourceLocation, n string, k ScopeKind) *VariableNode {
	node := &VariableNode{
		Base: Base{Location: sl},
		Name: n,
		Kind: k,
	}
	return node
}

type LambdaNode struct {
	Base
	VarList []*VariableNode
	Expr    ExprNode
}

func NewLambdaNode(sl file.SourceLocation, vl []*VariableNode, e ExprNode) *LambdaNode {
	node := &LambdaNode{
		Base:    Base{Location: sl},
		VarList: vl,
		Expr:    e,
	}
	return node
}

type LetrecVarExprItem struct {
	Variable *VariableNode
	Expr     ExprNode
}

type LetrecNode struct {
	Base
	VarExprList []*LetrecVarExprItem
	Expr        ExprNode
}

func NewLetrecNode(sl file.SourceLocation, vel []*LetrecVarExprItem, e ExprNode) *LetrecNode {
	node := &LetrecNode{
		Base:        Base{Location: sl},
		VarExprList: vel,
		Expr:        e,
	}
	return node
}

type IfNode struct {
	Base
	Cond    ExprNode
	Branch1 ExprNode
	Branch2 ExprNode
}

func NewIfNode(sl file.SourceLocation, c, b1, b2 ExprNode) *IfNode {
	node := &IfNode{
		Base:    Base{Location: sl},
		Cond:    c,
		Branch1: b1,
		Branch2: b2,
	}
	return node
}

type CallNode struct {
	Base
	Callee  ExprNode
	ArgList []ExprNode
}

func NewCallNode(sl file.SourceLocation, c ExprNode, al []ExprNode) *CallNode {
	node := &CallNode{
		Base:    Base{Location: sl},
		Callee:  c,
		ArgList: al,
	}
	return node
}

type SequenceNode struct {
	Base
	ExprList []ExprNode
}

func NewSequenceNode(sl file.SourceLocation, el []ExprNode) *SequenceNode {
	node := &SequenceNode{
		Base:     Base{Location: sl},
		ExprList: el,
	}
	return node
}

type QueryNode struct {
	Base
	Variable *VariableNode
	ExprBox  []ExprNode
}

func NewQueryNode(sl file.SourceLocation, v *VariableNode, eb []ExprNode) *QueryNode {
	node := &QueryNode{
		Base:     Base{Location: sl},
		Variable: v,
		ExprBox:  eb,
	}
	return node
}

type AccessNode struct {
	Base
	Variable *VariableNode
	Expr     ExprNode
}

func NewAccessNode(sl file.SourceLocation, v *VariableNode, e ExprNode) *AccessNode {
	node := &AccessNode{
		Base:     Base{Location: sl},
		Variable: v,
		Expr:     e,
	}
	return node
}
