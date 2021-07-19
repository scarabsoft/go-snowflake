package internal

import (
	"errors"
	"fmt"
	"io"
	"time"
)

var (
	ErrMaxSequenceOutOfRange = errors.New(fmt.Sprintf("maxSequence is capped to %d", MaxSequence))
)

// Sequence is contains information about the current sequence consisting of millis and iteration or an error
type Sequence struct {
	// number of passed since epoch
	Millis uint64
	// current iteration within the same milliseconds
	// this increments if and only if multiple sequences gets generated at the same millis
	// it can be max 16383
	Iteration uint16
	// error which occurred during the generation of the sequence
	Error error
}

func sequenceOk(millis uint64, it uint16) Sequence {
	return Sequence{Millis: millis, Iteration: it}
}

func sequenceError(err error) Sequence {
	return Sequence{Error: err}
}

// SequenceProvider generates sequences which will be used by the snowflake generator
type SequenceProvider interface {
	io.Closer

	//Generates the next sequence which must be unique, otherwise it will result in duplicated IDs
	Sequence() <-chan Sequence
}

type sequenceProviderImpl struct {
	clock        Clock
	maxIteration uint16

	closeChan   chan struct{}
	requestChan chan chan Sequence

	currentMillis    uint64
	currentIteration uint16
}

func (s *sequenceProviderImpl) Close() error {
	s.closeChan <- struct{}{}
	return nil
}

func (s *sequenceProviderImpl) Sequence() <-chan Sequence {
	r := make(chan Sequence)
	go func() {
		<- r
	}()
	return r
}

func (s *sequenceProviderImpl) generateNextSequence() Sequence {
	millis, err := s.clock.Millis()
	if err != nil {
		return sequenceError(err)
	}

	if millis < s.currentMillis {
		return sequenceError(ErrClockNotMonotonic)
	}

	if millis != s.currentMillis {
		s.currentMillis = millis
		s.currentIteration = 0
	}

	if s.currentIteration >= s.maxIteration {
		time.Sleep(1 * time.Millisecond)
		return s.generateNextSequence()
	}

	s.currentIteration += 1
	return sequenceOk(s.currentMillis, s.currentIteration)
}

func (s *sequenceProviderImpl) run() {
	go func() {
		for {
			select {
			case <-s.closeChan:
				return
			case responseChannel := <-s.requestChan:
				responseChannel <- s.generateNextSequence()
			}
		}
	}()
}

// Returns and starts a new sequence provider, can be stopped by invoking Close()
func NewSequenceProvider(clock Clock, maxSequence uint16) (*sequenceProviderImpl, error) {

	if maxSequence > MaxSequence {
		return nil, ErrMaxSequenceOutOfRange
	}

	r := &sequenceProviderImpl{
		clock:        clock,
		closeChan:    make(chan struct{}),
		requestChan:  make(chan chan Sequence, maxSequence),
		maxIteration: maxSequence,
	}
	r.run()
	return r, nil
}
