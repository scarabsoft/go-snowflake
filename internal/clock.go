package internal

import (
	"time"
)

type Clock interface {
	Seconds() uint64
}

type unixClockImpl struct {
	customEpoch uint64
}

func (u unixClockImpl) Seconds() uint64 {
	return uint64(time.Now().Unix()) - u.customEpoch
}

func NewUnixClockWithEpoch(epoch uint64) Clock {
	return &unixClockImpl{customEpoch: epoch}
}
