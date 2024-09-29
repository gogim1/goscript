package parser_test

import (
	"reflect"
	"testing"

	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
	. "github.com/gogim1/goscript/parser"
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
