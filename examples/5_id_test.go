package examples

import (
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"github.com/scarabsoft/go-snowflake"
	"testing"
)

type fakeClockImpl struct {
	value uint64
}

func (f fakeClockImpl) Seconds() uint64 {
	return f.value
}

func TestID(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	gen, err := snowflake.NewGenerator(
		snowflake.WithClock(fakeClockImpl{value: 1647619145}),
		snowflake.WithNodeID(128),
	)
	assert.That(err, is.Nil())

	testInstance, err := gen.Next()
	assert.That(err, is.Nil())
	assert.That(testInstance.NodeID(), is.EqualTo(uint8(128)))

	assert.That(testInstance.Iteration(), is.EqualTo(uint16(1)))

	assert.That(testInstance.Seconds(), is.EqualTo(uint64(1647619145)))
	assert.That(testInstance.Minutes(), is.EqualTo(uint64(27460319)))
	assert.That(testInstance.Hours(), is.EqualTo(uint64(457671)))
	assert.That(testInstance.Days(), is.EqualTo(uint64(19069)))
	assert.That(testInstance.Weeks(), is.EqualTo(uint64(2724)))

	assert.That(testInstance.String(), is.EqualTo("6910615572447233"))
}

func TestFrom(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	testInstance := snowflake.From(6823236456859828225)
	assert.That(testInstance.ID(), is.EqualTo(uint64(6823236456859828225)))

}
