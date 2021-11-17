package internal

import (
	"errors"
	"fmt"
)

var (
	ErrClockNotMonotonic     = errors.New("clock is not monotonic")
	ErrMaxSequenceOutOfRange = errors.New(fmt.Sprintf("maxSequence is capped to %d", MaxSequence))
)
