package internal

type SnowflakeGenerator interface {
	Next() <-chan Result
}

type snowFlakeGeneratorImpl struct {
	seqProvider SequenceProvider
	nodeID      uint8
}

func (s *snowFlakeGeneratorImpl) Next() <-chan Result {
	r := make(chan Result)
	go func() {
		defer close(r)
		seq := <-s.seqProvider.Sequence()
		if seq.Error != nil {
			r <- Result{0, seq.Error}
			return
		}

		id := seq.Millis << uint64(totalBits-epochBits)
		id |= uint64(s.nodeID) << uint64(totalBits-epochBits-nodeBits)
		id |= uint64(seq.Iteration)

		r <- Result{id, nil}
	}()
	return r
}

func NewGenerator(seq SequenceProvider, node NodeProvider) (SnowflakeGenerator, error) {
	nodeID, err := node.ID()
	if err != nil {
		return nil, err
	}
	return &snowFlakeGeneratorImpl{
		seqProvider: seq,
		nodeID:      nodeID,
	}, nil
}
