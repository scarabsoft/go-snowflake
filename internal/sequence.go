package internal

import (
	"sync"
	"time"
)

// Sequence is contains information about the current sequence consisting of millis and iteration or an error
type Sequence struct {
	// number of passed since epoch
	Seconds uint64
	// current iteration within the same milliseconds
	// this increments if and only if multiple sequences gets generated at the same millis
	// it can be max 16383
	Iteration uint16
	// error which occurred during the generation of the sequence
	Error error
}

func sequenceOk(seconds uint64, it uint16) Sequence {
	return Sequence{Seconds: seconds, Iteration: it}
}

func sequenceError(err error) Sequence {
	return Sequence{Error: err}
}

// SequenceProvider generates sequences which will be used by the snowflake generator
type SequenceProvider interface {
	// Sequence generates the next sequence which must be unique, otherwise it will result in duplicated IDs
	Sequence() Sequence
}

type sequenceProviderImpl struct {
	clock        Clock
	maxIteration uint16
	lock         sync.Mutex

	currentSeconds   uint64
	currentIteration uint16
}

func (s *sequenceProviderImpl) Sequence() Sequence {
	s.lock.Lock()

	secondsSinceEpoch := s.clock.Seconds()

	if secondsSinceEpoch < s.currentSeconds {
		s.lock.Unlock()
		return sequenceError(ErrClockNotMonotonic)
	}

	if secondsSinceEpoch != s.currentSeconds {
		s.currentSeconds = secondsSinceEpoch
		s.currentIteration = 0
	}

	if s.currentIteration >= s.maxIteration {
		time.Sleep(250 * time.Millisecond)
		s.lock.Unlock()
		return s.Sequence()
	}

	s.currentIteration += 1
	defer s.lock.Unlock()
	return sequenceOk(s.currentSeconds, s.currentIteration)
}

//NewSequenceProvider returns and starts a new sequence provider, can be stopped by invoking Close()
func NewSequenceProvider(clock Clock, maxSequence uint16) (*sequenceProviderImpl, error) {

	if maxSequence > MaxSequence {
		return nil, ErrMaxSequenceOutOfRange
	}

	r := &sequenceProviderImpl{
		clock:        clock,
		lock:         sync.Mutex{},
		maxIteration: maxSequence,
	}

	return r, nil
}
