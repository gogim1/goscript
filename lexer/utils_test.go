package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexUtils(t *testing.T) {
	t.Run("count_trailing_escape", func(t *testing.T) {
		assert.Equal(t, 0, countTrailingEscape([]rune(``)))
		assert.Equal(t, 1, countTrailingEscape([]rune(`\`)))
		assert.Equal(t, 3, countTrailingEscape([]rune(`\\\`)))
		assert.Equal(t, 0, countTrailingEscape([]rune(`123`)))
		assert.Equal(t, 0, countTrailingEscape([]rune(`你好`)))
	})
}
