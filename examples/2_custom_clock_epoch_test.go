package examples

import (
	"fmt"
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"github.com/scarabsoft/go-snowflake"
	"testing"
)

type incrementClockImpl struct {
	value uint64
}

func (i incrementClockImpl) Seconds() uint64 {
	i.value++
	return i.value
}

func TestCustomClockEpoch(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	gen, err := snowflake.NewGenerator(
		snowflake.WithClock(incrementClockImpl{}),
		snowflake.WithNodeID(1),
	)

	assert.That(err, is.Nil())

	r, err := gen.Next()
	assert.That(err, is.Nil())
	binId := fmt.Sprintf("%064b", r.ID())
	assert.That(binId, is.EqualTo("0000000000000000000000000000000000000000010000000100000000000001"))

}
