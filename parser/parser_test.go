package parser_test

import (
	"reflect"
	"testing"

	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
	. "github.com/gogim1/goscript/parser"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		node  ExprNode
	}{
		{
			"+1",
			&NumberNode{
				Base:        Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				Numerator:   1,
				Denominator: 1,
			},
		},
		{
			"-3/6",
			&NumberNode{
				Base:        Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				Numerator:   -1,
				Denominator: 2,
			},
		},
		{
			"123.45",
			&NumberNode{
				Base:        Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				Numerator:   2469,
				Denominator: 20,
			},
		},
		{
			`"123.45"`,
			&StringNode{
				Base:  Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				Value: file.NewSource(`123.45`),
			},
		},
		{
			`"\\\"\t\n"`,
			&StringNode{
				Base:  Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				Value: file.NewSource("\\\"\t\n"),
			},
		},
		{
			`a`,
			&VariableNode{
				Base: Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				Name: file.NewSource("a"),
			},
		},
		{
			`lambda () { 1 }`,
			&LambdaNode{
				Base:    Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				VarList: []*VariableNode{},
				Expr: &NumberNode{
					Base:        Base{Location: file.SourceLocation{Line: 1, Col: 13}},
					Numerator:   1,
					Denominator: 1,
				},
			},
		},
		{
			`lambda (a b) { c }`,
			&LambdaNode{
				Base: Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				VarList: []*VariableNode{
					{
						Base: Base{Location: file.SourceLocation{Line: 1, Col: 9}},
						Name: file.NewSource("a"),
					},
					{
						Base: Base{Location: file.SourceLocation{Line: 1, Col: 11}},
						Name: file.NewSource("b"),
					},
				},
				Expr: &VariableNode{
					Base: Base{Location: file.SourceLocation{Line: 1, Col: 16}},
					Name: file.NewSource("c"),
				},
			},
		},
		{
			`letrec () { 1/42 }`,
			&LetrecNode{
				Base:        Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				VarExprList: []*LetrecVarExprItem{},
				Expr: &NumberNode{
					Base:        Base{Location: file.SourceLocation{Line: 1, Col: 13}},
					Numerator:   1,
					Denominator: 42,
				},
			},
		},
		{
			`letrec (a=1 b=a) { "str" }`,
			&LetrecNode{
				Base: Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				VarExprList: []*LetrecVarExprItem{
					{
						Variable: &VariableNode{
							Base: Base{Location: file.SourceLocation{Line: 1, Col: 9}},
							Name: file.NewSource("a"),
						},
						Expr: &NumberNode{
							Base:        Base{Location: file.SourceLocation{Line: 1, Col: 11}},
							Numerator:   1,
							Denominator: 1,
						},
					},
					{
						Variable: &VariableNode{
							Base: Base{Location: file.SourceLocation{Line: 1, Col: 13}},
							Name: file.NewSource("b"),
						},
						Expr: &VariableNode{
							Base: Base{Location: file.SourceLocation{Line: 1, Col: 15}},
							Name: file.NewSource("a"),
						},
					},
				},
				Expr: &StringNode{
					Base:  Base{Location: file.SourceLocation{Line: 1, Col: 20}},
					Value: file.NewSource("str"),
				},
			},
		},
		{
			`if 1 then 2 else b`,
			&IfNode{
				Base: Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				Cond: &NumberNode{
					Base:        Base{Location: file.SourceLocation{Line: 1, Col: 4}},
					Numerator:   1,
					Denominator: 1,
				},
				Branch1: &NumberNode{
					Base:        Base{Location: file.SourceLocation{Line: 1, Col: 11}},
					Numerator:   2,
					Denominator: 1,
				},
				Branch2: &VariableNode{
					Base: Base{Location: file.SourceLocation{Line: 1, Col: 18}},
					Name: file.NewSource("b"),
				},
			},
		},
		{
			`(put "hello world")`,
			&CallNode{
				Base: Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				Callee: &IntrinsicNode{
					Base: Base{Location: file.SourceLocation{Line: 1, Col: 2}},
					Name: file.NewSource("put"),
				},
				ArgList: []ExprNode{
					&StringNode{
						Base:  Base{Location: file.SourceLocation{Line: 1, Col: 6}},
						Value: file.NewSource("hello world"),
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tokens, err := lexer.Lex(file.NewSource(test.input))
			if err != nil {
				t.Errorf("%s:\n%v", test.input, err)
				return
			}
			node, err := Parse(tokens)
			if err != nil {
				t.Errorf("%s:\n%v", test.input, err)
				return
			}
			if !reflect.DeepEqual(node, test.node) {
				t.Errorf("%s:\ngot\n\t%+v\nexpected\n\t%+v", test.input, node, test.node)
			}
		})
	}
}

func TestParse_error(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{
			"unsupported escape sequence",
			`"\b"`,
		},
		{
			"malformed lambda #1",
			`lambda () {}`,
		},
		{
			"malformed lambda #2",
			`lambda`,
		},
		{
			"malformed letrec #1",
			`letrec () {}`,
		},
		{
			"malformed letrec #1",
			`letrec (a) { 1 }`,
		},
		{
			"incorrect variable name",
			`letrec (exit) { 1 }`,
		},
		{
			"zero-length sequence",
			`[]`,
		},
		{
			"malformed call #1",
			`()`,
		},
		{
			"malformed call #2",
			`(`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, err := lexer.Lex(file.Source(test.input))
			assert.Nil(t, err)
			_, err = Parse(tokens)
			assert.NotNil(t, err)
			t.Log(err)
		})
	}
}
