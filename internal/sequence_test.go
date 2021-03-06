package internal

import (
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"testing"
	"time"
)

type fakeClock struct {
	value uint64
}

func (f fakeClock) Seconds() uint64 {
	return f.value
}

func TestSequenceOk(t *testing.T) {
	assert := hamcrest.NewAssertion(t)
	r := sequenceOk(10, 20)
	assert.That(r.Seconds, is.EqualTo(uint64(10)))
	assert.That(r.Iteration, is.EqualTo(uint16(20)))
	assert.That(r.Error, is.Nil())
}

func TestGenerateNextSequence(t *testing.T) {
	t.Run("clock skew", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		testInstance := &sequenceProviderImpl{currentSeconds: 10, clock: fakeClock{6}}
		seq := testInstance.Sequence()
		assert.That(seq.Seconds, is.EqualTo(uint64(0)))
		assert.That(seq.Iteration, is.EqualTo(uint16(0)))
		assert.That(seq.Error, is.EqualTo(ErrClockNotMonotonic))
	})
	t.Run("seq exhaustion", func(t *testing.T) {
		//wg := sync.WaitGroup{}
		c := make(chan struct{})
		testInstance := &sequenceProviderImpl{clock: fakeClock{10}, maxIteration: 0}

		go func() {
			// this is expected to block forever as the clock does not make any progress
			testInstance.Sequence()
			c <- struct{}{}
		}()

		for {
			select {
			case <-c:
				t.Fatal("this should never ever called")
			case <-time.After(1 * time.Millisecond):
				return
			}
		}
	})
	t.Run("ok", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		testInstance := &sequenceProviderImpl{clock: fakeClock{10}, maxIteration: 2}
		seq := testInstance.Sequence()
		assert.That(seq.Seconds, is.EqualTo(uint64(10)))
		assert.That(seq.Iteration, is.EqualTo(uint16(1)))
		assert.That(seq.Error, is.Nil())
	})
}

type incrementalClock struct {
	currentTime        uint64
	currentCallCounter int
}

func (i *incrementalClock) Seconds() uint64 {
	i.currentCallCounter++
	if i.currentCallCounter > 10 {
		i.currentCallCounter = 0
		i.currentTime++
	}

	return i.currentTime
}

func TestSequenceProviderImpl(t *testing.T) {
	assert := hamcrest.NewAssertion(t)
	testInstance, err := NewSequenceProvider(&incrementalClock{}, 100)
	assert.That(err, is.Nil())

	for j := 0; j < 2; j++ {
		for i := 0; i < 10; i++ {
			seq := testInstance.Sequence()
			assert.That(seq.Seconds, is.EqualTo(uint64(j)))
			assert.That(seq.Iteration, is.EqualTo(uint16(i+1)))
			assert.That(seq.Error, is.Nil())
		}
	}
}

func TestSequenceProvider_SequenceExhaustion(t *testing.T) {
	assert := hamcrest.NewAssertion(t)
	testInstance, err := NewSequenceProvider(&incrementalClock{}, 5)
	assert.That(err, is.Nil())

	for i := 0; i < 5; i++ {
		seq := testInstance.Sequence()
		assert.That(seq.Seconds, is.EqualTo(uint64(0)))
		assert.That(seq.Iteration, is.EqualTo(uint16(i+1)))
		assert.That(seq.Error, is.Nil())
	}

	for i := 0; i < 5; i++ {
		seq := testInstance.Sequence()
		assert.That(seq.Seconds, is.EqualTo(uint64(1)))
		assert.That(seq.Iteration, is.EqualTo(uint16(i+1)))
		assert.That(seq.Error, is.Nil())
	}
}
