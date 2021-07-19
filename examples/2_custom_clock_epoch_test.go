package examples

import (
	"fmt"
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"github.com/scarabsoft/go-snowflake"
	"sync"
	"testing"
	"time"
)

type incrementClockImpl struct {
	value uint64
	lock  sync.Mutex
}

func (i *incrementClockImpl) Millis() (uint64, error) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.value++
	return i.value, nil
}

func TestCustomClockEpoch(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	now := uint64(time.Now().UnixNano())
	gen, err := snowflake.New(
		snowflake.WithClock(snowflake.NewUnixClockWithEpoch(now)),
	)

	assert.That(err, is.Nil())

	r := <-gen.Next()
	binId := fmt.Sprintf("%064b", r.ID)
	assert.That(binId, is.EqualTo("0000000000000000000000000000000000000000000000000100000000000001"))

}
