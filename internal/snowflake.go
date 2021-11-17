package internal

type SnowflakeGenerator interface {
	Next() (uint64, error)
}

type snowFlakeGeneratorImpl struct {
	seqProvider SequenceProvider
	nodeID      uint8
}

func (s *snowFlakeGeneratorImpl) Next() (uint64, error) {
	seq := s.seqProvider.Sequence()
	if seq.Error != nil {
		return 0, seq.Error
	}

	id := seq.Seconds << uint64(totalBits-epochBits)
	id |= uint64(s.nodeID) << uint64(totalBits-epochBits-nodeBits)
	id |= uint64(seq.Iteration)

	return id, nil
}

func NewGenerator(seq SequenceProvider, node NodeIDProvider) (SnowflakeGenerator, error) {
	return &snowFlakeGeneratorImpl{
		seqProvider: seq,
		nodeID:      node.ID(),
	}, nil
}
