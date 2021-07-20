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

func (f fakeClockImpl) Millis() (uint64, error) {
	return f.value, nil
}

func TestID(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	gen, err := snowflake.New(
		snowflake.WithClock(fakeClockImpl{value: 1626786340918}),
		snowflake.WithNodeID(128),
	)
	assert.That(err, is.Nil())

	testInstance, err := gen.Next()
	assert.That(err, is.Nil())
	assert.That(testInstance.NodeID(), is.EqualTo(uint8(128)))

	assert.That(testInstance.Iteration(), is.EqualTo(uint16(1)))

	assert.That(testInstance.Millis(), is.EqualTo(uint64(1626786340918)))
	assert.That(testInstance.Seconds(), is.EqualTo(uint64(1626786340)))
	assert.That(testInstance.Minutes(), is.EqualTo(uint64(27113105)))
	assert.That(testInstance.Hours(), is.EqualTo(uint64(451885)))
	assert.That(testInstance.Days(), is.EqualTo(uint64(18828)))
	assert.That(testInstance.Weeks(), is.EqualTo(uint64(2689)))

	assert.That(testInstance.String(), is.EqualTo("6823236456859828225"))
}

func TestFrom(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	testInstance := snowflake.From(6823236456859828225)
	assert.That(testInstance.ID(), is.EqualTo(uint64(6823236456859828225)))

}
