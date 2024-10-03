package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexUtils(t *testing.T) {
	t.Run("count_trailing_escape", func(t *testing.T) {
		assert.Equal(t, 0, countTrailingEscape(``))
		assert.Equal(t, 1, countTrailingEscape(`\`))
		assert.Equal(t, 3, countTrailingEscape(`\\\`))
		assert.Equal(t, 0, countTrailingEscape(`123`))
		assert.Equal(t, 0, countTrailingEscape(`你好`))
	})
}
