package internal

import (
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"testing"
	"time"
)

func TestErrClockNotMonotonic(t *testing.T) {
	assert := hamcrest.NewAssertion(t)
	assert.That(ErrClockNotMonotonic.Error(), is.EqualTo("clock is not monotonic"))
}

func TestNewUnixClockWithEpoch(t *testing.T) {
	t.Run("provided epoch is 0", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := NewUnixClockWithEpoch(0)
		r := testInstance.Seconds()
		assert.That(r, is.EqualTo(uint64(time.Now().Unix())))
	})

	t.Run("provided epoch which is now", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := NewUnixClockWithEpoch(uint64(time.Now().Unix()))
		r := testInstance.Seconds()
		assert.That(r, is.EqualTo(uint64(0)))
	})
}
