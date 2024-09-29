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

type Base struct {
	Location file.SourceLocation
}

func (n *Base) GetLocation() file.SourceLocation {
	return n.Location
}

func (n *Base) SetLocation(sl file.SourceLocation) {
	n.Location = sl
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
	Value []rune
}

func NewStringNode(sl file.SourceLocation, v []rune) *StringNode {
	node := &StringNode{
		Base:  Base{Location: sl},
		Value: v,
	}
	return node
}

type IntrinsicNode struct {
	Base
	Name []rune
}

func NewIntrinsicNode(sl file.SourceLocation, n []rune) *IntrinsicNode {
	node := &IntrinsicNode{
		Base: Base{Location: sl},
		Name: n,
	}
	return node
}

type VariableNode struct {
	Base
	Name []rune
}

func (n *VariableNode) IsLex() bool {
	return unicode.IsLower(n.Name[0])
}

func (n *VariableNode) IsDyn() bool {
	return unicode.IsUpper(n.Name[0])
}

func NewVariableNode(sl file.SourceLocation, n []rune) *VariableNode {
	node := &VariableNode{
		Base: Base{Location: sl},
		Name: n,
	}
	return node
}

type LambdaNode struct {
	Base
	VarList []ExprNode
	Expr    ExprNode
}

func NewLambdaNode(sl file.SourceLocation, vl []ExprNode, e ExprNode) *LambdaNode {
	node := &LambdaNode{
		Base:    Base{Location: sl},
		VarList: vl,
		Expr:    e,
	}
	return node
}

type LetrecVarExprItem struct {
	Variable VariableNode
	Expr     ExprNode
}

func NewLetrecVarExprItem(v VariableNode, e ExprNode) *LetrecVarExprItem {
	node := &LetrecVarExprItem{
		Variable: v,
		Expr:     e,
	}
	return node
}

type LetrecNode struct {
	Base
	VarExprList []LetrecVarExprItem
	Expr        ExprNode
}

func NewLetrecNode(sl file.SourceLocation, vel []LetrecVarExprItem, e ExprNode) *LetrecNode {
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
	Variable VariableNode
	ExprBox  []ExprNode
}

func NewQueryNode(sl file.SourceLocation, v VariableNode, eb []ExprNode) *QueryNode {
	node := &QueryNode{
		Base:     Base{Location: sl},
		Variable: v,
		ExprBox:  eb,
	}
	return node
}

type AccessNode struct {
	Base
	Variable VariableNode
	Expr     ExprNode
}

func NewAccessNode(sl file.SourceLocation, v VariableNode, e ExprNode) *AccessNode {
	node := &AccessNode{
		Base:     Base{Location: sl},
		Variable: v,
		Expr:     e,
	}
	return node
}
