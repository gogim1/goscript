package runtime_test

import (
	"testing"

	"github.com/gogim1/goscript/file"
	"github.com/gogim1/goscript/lexer"
	"github.com/gogim1/goscript/parser"
	. "github.com/gogim1/goscript/runtime"
)

func TestRuntime(t *testing.T) {
	tests := []struct {
		input, value string
	}{
		{`10/5`, `2`},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tokens, err := lexer.Lex(file.NewSource(test.input))
			if err != nil {
				t.Errorf("%s:\n%v", test.input, err)
				return
			}
			node, err := parser.Parse(tokens)
			if err != nil {
				t.Errorf("%s:\n%v", test.input, err)
				return
			}
			state := NewState(node)
			if err := state.Execute(); err != nil {
				t.Errorf("%s:\n%v", test.input, err)
				return
			}
			if state.Value() != test.value {
				t.Errorf("%s:\ngot\n\t%+v\nexpected\n\t%+v", test.input, state.Value(), test.value)
			}
		})
	}
}
