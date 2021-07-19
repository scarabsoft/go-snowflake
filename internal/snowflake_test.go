package internal

import (
	"fmt"
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"testing"
)

type countUpClock struct {
	value uint64
}

func (c *countUpClock) Millis() (uint64, error) {
	c.value++
	return c.value, nil
}

func TestSnowFlakeGeneratorImpl_Next(t *testing.T) {
	t.Run("within same ms", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		seqProvider, err := NewSequenceProvider(fakeClock{10}, 10)
		assert.That(err, is.Nil())

		testInstance, err := NewGenerator(
			seqProvider,
			fixedNodeProviderImpl{42},
		)
		assert.That(err, is.Nil())

		t.Run("1", func(t *testing.T) {
			assert := hamcrest.NewAssertion(t)

			r := testInstance.Next()
			assert.That(r.Error, is.Nil())
			bin := fmt.Sprintf("%064b", r.ID)
			assert.That(bin, is.EqualTo("0000000000000000000000000000000000000010100010101000000000000001"))
		})

		t.Run("2", func(t *testing.T) {
			assert := hamcrest.NewAssertion(t)

			r := testInstance.Next()
			assert.That(r.Error, is.Nil())
			bin := fmt.Sprintf("%064b", r.ID)
			assert.That(bin, is.EqualTo("0000000000000000000000000000000000000010100010101000000000000010"))
		})

		t.Run("3", func(t *testing.T) {
			assert := hamcrest.NewAssertion(t)

			r := testInstance.Next()
			assert.That(r.Error, is.Nil())
			bin := fmt.Sprintf("%064b", r.ID)
			assert.That(bin, is.EqualTo("0000000000000000000000000000000000000010100010101000000000000011"))
		})
	})

	t.Run("once per ms", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		seqProvider, err := NewSequenceProvider(&countUpClock{}, 10)
		assert.That(err, is.Nil())

		testInstance, err := NewGenerator(
			seqProvider,
			fixedNodeProviderImpl{42},
		)
		assert.That(err, is.Nil())

		t.Run("1", func(t *testing.T) {
			assert := hamcrest.NewAssertion(t)

			r := testInstance.Next()
			assert.That(r.Error, is.Nil())
			bin := fmt.Sprintf("%064b", r.ID)
			assert.That(bin, is.EqualTo("0000000000000000000000000000000000000000010010101000000000000001"))
		})

		t.Run("2", func(t *testing.T) {
			assert := hamcrest.NewAssertion(t)

			r := testInstance.Next()
			assert.That(r.Error, is.Nil())
			bin := fmt.Sprintf("%064b", r.ID)
			assert.That(bin, is.EqualTo("0000000000000000000000000000000000000000100010101000000000000001"))
		})

		t.Run("3", func(t *testing.T) {
			assert := hamcrest.NewAssertion(t)

			r := testInstance.Next()
			assert.That(r.Error, is.Nil())
			bin := fmt.Sprintf("%064b", r.ID)
			assert.That(bin, is.EqualTo("0000000000000000000000000000000000000000110010101000000000000001"))
		})
	})

	t.Run("different host", func(t *testing.T) {

		t.Run("1", func(t *testing.T) {
			assert := hamcrest.NewAssertion(t)

			seqProvider, err := NewSequenceProvider(fakeClock{10}, 10)
			assert.That(err, is.Nil())

			testInstance, err := NewGenerator(
				seqProvider,
				fixedNodeProviderImpl{1},
			)
			assert.That(err, is.Nil())

			r := testInstance.Next()
			assert.That(r.Error, is.Nil())
			bin := fmt.Sprintf("%064b", r.ID)
			assert.That(bin, is.EqualTo("0000000000000000000000000000000000000010100000000100000000000001"))
		})

		t.Run("2", func(t *testing.T) {
			assert := hamcrest.NewAssertion(t)

			seqProvider, err := NewSequenceProvider(fakeClock{10}, 10)
			assert.That(err, is.Nil())

			testInstance, err := NewGenerator(
				seqProvider,
				fixedNodeProviderImpl{2},
			)
			assert.That(err, is.Nil())

			r := testInstance.Next()
			assert.That(r.Error, is.Nil())
			bin := fmt.Sprintf("%064b", r.ID)
			assert.That(bin, is.EqualTo("0000000000000000000000000000000000000010100000001000000000000001"))
		})

	})
}
