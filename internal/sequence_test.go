package internal

import (
	"errors"
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"sync"
	"testing"
	"time"
)

var givenErrClock = errors.New("givenErrClock")

type errorClock struct{}

func (e errorClock) Millis() (uint64, error) {
	return 0, givenErrClock
}

type fakeClock struct {
	value uint64
}

func (f fakeClock) Millis() (uint64, error) {
	return f.value, nil
}

func TestSequenceOk(t *testing.T) {
	assert := hamcrest.NewAssertion(t)
	r := sequenceOk(10, 20)
	assert.That(r.Millis, is.EqualTo(uint64(10)))
	assert.That(r.Iteration, is.EqualTo(uint16(20)))
	assert.That(r.Error, is.Nil())
}

func TestSequenceError(t *testing.T) {
	assert := hamcrest.NewAssertion(t)
	givenError := errors.New("givenError")
	r := sequenceError(givenError)
	assert.That(r.Millis, is.EqualTo(uint64(0)))
	assert.That(r.Iteration, is.EqualTo(uint16(0)))
	assert.That(r.Error, is.EqualTo(givenError))
}

func TestSequenceProviderImpl_Close(t *testing.T) {
	wg := sync.WaitGroup{}
	assert := hamcrest.NewAssertion(t)

	closeChan := make(chan struct{})
	testInstance := &sequenceProviderImpl{closeChan: closeChan}
	go func() {
		wg.Add(1)
		for {
			select {
			case <-closeChan:
				wg.Done()
				return
			case <-time.After(10 * time.Millisecond):
				t.Fatal("Close channel not received within 10 ms")
			}
		}
	}()

	err := testInstance.Close()
	assert.That(err, is.Nil())

	wg.Wait()
}

func TestGenerateNextSequence(t *testing.T) {
	t.Run("clock returns error", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		testInstance := &sequenceProviderImpl{clock: errorClock{}}

		seq := testInstance.generateNextSequence()
		assert.That(seq.Millis, is.EqualTo(uint64(0)))
		assert.That(seq.Iteration, is.EqualTo(uint16(0)))
		assert.That(seq.Error, is.EqualTo(givenErrClock))
	})
	t.Run("clock skew", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)
		testInstance := &sequenceProviderImpl{currentMillis: 10, clock: fakeClock{6}}
		seq := testInstance.generateNextSequence()
		assert.That(seq.Millis, is.EqualTo(uint64(0)))
		assert.That(seq.Iteration, is.EqualTo(uint16(0)))
		assert.That(seq.Error, is.EqualTo(ErrClockNotMonotonic))
	})
	t.Run("seq exhaustion", func(t *testing.T) {
		//wg := sync.WaitGroup{}
		c := make(chan struct{})
		testInstance := &sequenceProviderImpl{clock: fakeClock{10}, maxIteration: 0}

		go func() {
			// this is expected to block forever as the clock does not make any progress
			testInstance.generateNextSequence()
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
		seq := testInstance.generateNextSequence()
		assert.That(seq.Millis, is.EqualTo(uint64(10)))
		assert.That(seq.Iteration, is.EqualTo(uint16(1)))
		assert.That(seq.Error, is.Nil())
	})
}

type incrementalClock struct {
	currentTime        uint64
	currentCallCounter int
}

func (i *incrementalClock) Millis() (uint64, error) {
	i.currentCallCounter++
	if i.currentCallCounter > 10 {
		i.currentCallCounter = 0
		i.currentTime++
	}

	return i.currentTime, nil
}

func TestSequenceProviderImpl(t *testing.T) {
	assert := hamcrest.NewAssertion(t)
	testInstance, err := NewSequenceProvider(&incrementalClock{}, 100)
	assert.That(err, is.Nil())

	for j := 0; j < 2; j++ {
		for i := 0; i < 10; i++ {
			c := testInstance.Sequence()
			seq := <-c
			assert.That(seq.Millis, is.EqualTo(uint64(j)))
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
		c := testInstance.Sequence()
		seq := <-c
		assert.That(seq.Millis, is.EqualTo(uint64(0)))
		assert.That(seq.Iteration, is.EqualTo(uint16(i+1)))
		assert.That(seq.Error, is.Nil())
	}

	for i := 0; i < 5; i++ {
		c := testInstance.Sequence()
		seq := <-c
		assert.That(seq.Millis, is.EqualTo(uint64(1)))
		assert.That(seq.Iteration, is.EqualTo(uint16(i+1)))
		assert.That(seq.Error, is.Nil())
	}
}
