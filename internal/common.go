package internal

import "math"

const (
	totalBits    = 64
	epochBits    = 42
	nodeBits     = 8
	sequenceBits = 14
)

var (
	MaxSequence = (uint16)(math.Pow(2, sequenceBits) - 1)
)
