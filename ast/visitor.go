package ast

import "github.com/gogim1/goscript/file"

type Visitor interface {
	VisitNumberNode(*NumberNode) *file.Error
	VisitStringNode(*StringNode) *file.Error
	VisitIntrinsicNode(*IntrinsicNode) *file.Error
	VisitVariableNode(*VariableNode) *file.Error
	VisitLambdaNode(*LambdaNode) *file.Error
	VisitLetrecNode(*LetrecNode) *file.Error
	VisitIfNode(*IfNode) *file.Error
	VisitCallNode(*CallNode) *file.Error
	VisitSequenceNode(*SequenceNode) *file.Error
	VisitAccessNode(*AccessNode) *file.Error
}

func (n *NumberNode) Accept(v Visitor) *file.Error {
	return v.VisitNumberNode(n)
}

func (n *StringNode) Accept(v Visitor) *file.Error {
	return v.VisitStringNode(n)
}

func (n *IntrinsicNode) Accept(v Visitor) *file.Error {
	return v.VisitIntrinsicNode(n)
}

func (n *VariableNode) Accept(v Visitor) *file.Error {
	return v.VisitVariableNode(n)
}

func (n *LambdaNode) Accept(v Visitor) *file.Error {
	return v.VisitLambdaNode(n)
}

func (n *LetrecNode) Accept(v Visitor) *file.Error {
	return v.VisitLetrecNode(n)
}

func (n *IfNode) Accept(v Visitor) *file.Error {
	return v.VisitIfNode(n)
}

func (n *CallNode) Accept(v Visitor) *file.Error {
	return v.VisitCallNode(n)
}

func (n *SequenceNode) Accept(v Visitor) *file.Error {
	return v.VisitSequenceNode(n)
}

func (n *AccessNode) Accept(v Visitor) *file.Error {
	return v.VisitAccessNode(n)
}
