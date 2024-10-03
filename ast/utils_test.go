package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUtils(t *testing.T) {
	t.Run("gcd", func(t *testing.T) {
		assert.Equal(t, 2, gcd(0, 2))
		assert.Equal(t, 1, gcd(2, 1))
		assert.Equal(t, 2, gcd(4, 6))
		assert.Equal(t, 2, gcd(-4, 6))
		assert.Equal(t, -2, gcd(-4, -6))
		assert.Equal(t, -2, gcd(4, -6))
	})
}
