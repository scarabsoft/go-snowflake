package internal

import (
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"testing"
)

func TestConstants(t *testing.T) {
	t.Run("totalBits", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		assert.That(totalBits, is.EqualTo(64))
	})
	t.Run("epochBits", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		assert.That(epochBits, is.EqualTo(42))
	})
	t.Run("nodeBits", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		assert.That(nodeBits, is.EqualTo(8))
	})
	t.Run("sequenceBits", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		assert.That(sequenceBits, is.EqualTo(14))
	})
}

func TestVariables(t *testing.T){
	t.Run("MaxSequence", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		assert.That(MaxSequence, is.EqualTo(uint16(16383)))
	})
}