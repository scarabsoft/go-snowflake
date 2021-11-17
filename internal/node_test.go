package internal

import (
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	assert := hamcrest.NewAssertion(t)
	testInstance := NewFixedNodeIdProvider(128)
	assert.That(testInstance.id, is.EqualTo(uint8(128)))
}

func TestFixedNodeProvider_ID(t *testing.T) {
	assert := hamcrest.NewAssertion(t)
	testInstance := NewFixedNodeIdProvider(128)
	r := testInstance.ID()
	assert.That(r, is.EqualTo(uint8(128)))
}
