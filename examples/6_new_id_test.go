package examples

import (
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"github.com/scarabsoft/go-snowflake"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestNewID(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	_ = os.Setenv("NODE_ID", "23")
	_ = os.Setenv("GENESIS_EPOCH_SECONDS", strconv.FormatInt(time.Now().Unix(), 10))

	result := snowflake.MustNewID()
	assert.That(result.ID(), is.EqualTo(uint64(376833)))

}
