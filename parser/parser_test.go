package parser_test

import (
	"testing"

	. "github.com/gogim1/goscript/ast"
	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
	. "github.com/gogim1/goscript/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				Value: `123.45`,
			},
		},
		{
			`"\\\"\t\n"`,
			&StringNode{
				Base:  Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				Value: "\\\"\t\n",
			},
		},
		{
			`a`,
			&VariableNode{
				Base: Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				Name: "a",
				Kind: Lexical,
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
			`lambda (a B) { c }`,
			&LambdaNode{
				Base: Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				VarList: []*VariableNode{
					{
						Base: Base{Location: file.SourceLocation{Line: 1, Col: 9}},
						Name: "a",
						Kind: Lexical,
					},
					{
						Base: Base{Location: file.SourceLocation{Line: 1, Col: 11}},
						Name: "B",
						Kind: Dynamic,
					},
				},
				Expr: &VariableNode{
					Base: Base{Location: file.SourceLocation{Line: 1, Col: 16}},
					Name: "c",
					Kind: Lexical,
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
							Name: "a",
							Kind: Lexical,
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
							Name: "b",
							Kind: Lexical,
						},
						Expr: &VariableNode{
							Base: Base{Location: file.SourceLocation{Line: 1, Col: 15}},
							Name: "a",
							Kind: Lexical,
						},
					},
				},
				Expr: &StringNode{
					Base:  Base{Location: file.SourceLocation{Line: 1, Col: 20}},
					Value: "str",
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
					Name: "b",
					Kind: Lexical,
				},
			},
		},
		{
			`(put "hello world")`,
			&CallNode{
				Base: Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				Callee: &IntrinsicNode{
					Base: Base{Location: file.SourceLocation{Line: 1, Col: 2}},
					Name: "put",
				},
				ArgList: []ExprNode{
					&StringNode{
						Base:  Base{Location: file.SourceLocation{Line: 1, Col: 6}},
						Value: "hello world",
					},
				},
			},
		},
		{
			`(user_defined_function)`,
			&CallNode{
				Base:    Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				ArgList: []ExprNode{},
				Callee: &VariableNode{
					Base: Base{Location: file.SourceLocation{Line: 1, Col: 2}},
					Name: "user_defined_function",
					Kind: Lexical,
				},
			},
		},
		{
			`&a b`,
			&AccessNode{
				Base: Base{Location: file.SourceLocation{Line: 1, Col: 1}},
				Variable: &VariableNode{
					Base: Base{Location: file.SourceLocation{Line: 1, Col: 2}},
					Name: "a",
					Kind: Lexical,
				},
				Expr: &VariableNode{
					Base: Base{Location: file.SourceLocation{Line: 1, Col: 4}},
					Name: "b",
					Kind: Lexical,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tokens, err := lexer.Lex(file.NewSource(test.input))
			require.Nil(t, err)

			node, err := Parse(tokens)
			assert.Nil(t, err)

			assert.Equal(t, test.node, node)
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
		{
			"incorrect access",
			`&letrec`,
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
