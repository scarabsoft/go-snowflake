package internal

import (
	"errors"
	"time"
)

var ErrClockNotMonotonic = errors.New("clock is not monotonic")

type Clock interface {
	Millis() (uint64, error)
}

type unixClockImpl struct {
	customEpoch uint64
}

func (u unixClockImpl) Millis() (uint64, error) {
	return (uint64(time.Now().UnixNano() ) - u.customEpoch) / 1e6, nil
}

func NewUnixClockWithEpoch(epoch uint64) Clock {
	return &unixClockImpl{customEpoch: epoch}
}
