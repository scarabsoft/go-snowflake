package examples

import (
	"fmt"
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/has"
	"github.com/scarabsoft/go-hamcrest/is"
	"github.com/scarabsoft/go-snowflake"
	"testing"
)

type fakeClock struct {
	value uint64
}

func (f fakeClock) Seconds() uint64 {
	return f.value
}

func TestCustomNodeId(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	clock := fakeClock{1337}

	nodeGen1, err := snowflake.NewGenerator(
		snowflake.WithNodeID(1),
		snowflake.WithClock(clock),
	)

	assert.That(err, is.Nil())

	nodeGen2, err := snowflake.NewGenerator(
		snowflake.WithNodeID(2),
		snowflake.WithClock(clock),
	)

	assert.That(err, is.Nil())

	r1, _ := nodeGen1.Next()
	r2, _ := nodeGen2.Next()

	assert.That(r2.ID(), is.GreaterThan(r1.ID()))

	binId1 := fmt.Sprintf("%064b", r1.ID())
	binId2 := fmt.Sprintf("%064b", r2.ID())

	assert.That(binId1, has.Prefix("00000000000000000000000000000001"))
	assert.That(binId2, has.Prefix("00000000000000000000000000000001"))

	assert.That(binId1, has.Suffix("00000000000001"))
	assert.That(binId2, has.Suffix("00000000000001"))

}
