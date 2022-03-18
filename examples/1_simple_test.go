package examples

import (
	"fmt"
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"github.com/scarabsoft/go-snowflake"
	"testing"
)

func TestSimple(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	gen, err := snowflake.NewGenerator()
	assert.That(err, is.Nil())

	var prev uint64 = 0

	for i := 0; i < 10; i++ {
		r, err := gen.Next()
		assert.That(err, is.Nil())
		fmt.Printf("%064b\n", r.ID())
		assert.That(r.ID(), is.GreaterThan(uint64(0)))

		assert.That(r.ID(), is.GreaterThan(prev))
		prev = r.ID()
	}
}
