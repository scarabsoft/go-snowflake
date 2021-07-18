package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"time"
)

const (
	totalBits    = 64
	epochBits    = 42
	nodeBits     = 8
	sequenceBits = 12
)

var ErrClockNotMonotonic = errors.New("clock is not monotonic")

var (
	//maxSequence = (int16)(math.Pow(2, sequenceBits) - 1)
	maxSequence = int16(4)
	maxNodeId   = (int16)(math.Pow(2, nodeBits) - 1)
)

type Clock interface {
	// Return the number of milliseconds passed since a specific epoch
	Millis() (int64, error)
}

type UnixClock struct{}

func (u UnixClock) Millis() (int64, error) {
	return time.Now().UnixNano() / 1e6, nil
}

type NodeProvider interface {
	ID() (int8, error)
}

type fixedNodeProvider struct {
	id int8
}

func (f fixedNodeProvider) ID() (int8, error) {
	return f.id, nil
}

type Result struct {
	ID    int64
	Error error
}

type Sequence struct {
	Millis    int64
	Iteration int16
	Error     error
}

func sequenceOk(millis int64, it int16) Sequence {
	return Sequence{Millis: millis, Iteration: it}
}

func sequenceError(err error) Sequence {
	return Sequence{Error: err}
}

type SequenceProvider interface {
	io.Closer
	Sequence() <-chan Sequence
}

type sequenceProviderImpl struct {
	clock        Clock
	maxIteration int16

	closeChan   chan struct{}
	requestChan chan chan Sequence

	currentMillis    int64
	currentIteration int16
}

func (s *sequenceProviderImpl) Close() error {
	s.closeChan <- struct{}{}
	return nil
}

func (s *sequenceProviderImpl) Sequence() <-chan Sequence {
	r := make(chan Sequence)
	go func() {
		s.requestChan <- r
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

func NewSequenceProvider(clock Clock, maxIteration int16) *sequenceProviderImpl {
	r := &sequenceProviderImpl{
		clock:        clock,
		closeChan:    make(chan struct{}),
		requestChan:  make(chan chan Sequence, maxSequence),
		maxIteration: maxIteration,
	}
	r.run()
	return r
}

type Generator interface {
	Next() <-chan Result
}
type generatorImpl struct {
	seqProvider  SequenceProvider
	nodeProvider NodeProvider
}

func (s *generatorImpl) Next() <-chan Result {
	r := make(chan Result)
	go func() {
		defer close(r)
		seq := <-s.seqProvider.Sequence()
		if seq.Error != nil {
			r <- Result{0, seq.Error}
			return
		}

		nodeId, err := s.nodeProvider.ID()
		if err != nil {
			r <- Result{0, err}
			return
		}

		id := seq.Millis << uint(totalBits-epochBits)
		id |= int64(nodeId << uint(totalBits-epochBits-nodeBits))
		id |= int64(seq.Iteration)

		r <- Result{id, nil}
	}()
	return r
}

//func (s *generatorImpl) Next() (int64, error) {
//	var result = millis, err := s.clock.Sequence()
//	if err != nil {
//		return 0, err
//	}
//
//}

//var sf = snowflake.NewSnowFlake()

func generateUniqueSequence(i int) {
	//seqID, err := sf.GenerateUniqueSequenceID()
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(i, "    ", "seqID::", seqID)
}

func main() {
	//FIXME document that
	fmt.Println("MaxSeq ", maxSequence)
	fmt.Println("MaxNodeId ", maxNodeId)

	c := UnixClock{}
	//fmt.Println(c.Sequence())

	provider := NewSequenceProvider(c, 1)

	for i := 0; i < 10; i++ {
		go func() {
			s := <-provider.Sequence()
			fmt.Println(s)
		}()

	}

	//time.Sleep(1 * time.Second)
	//
	//provider.Close()

	gen := generatorImpl{
		seqProvider:  provider,
		nodeProvider: fixedNodeProvider{2},
	}

	fmt.Println(<-gen.Next())
	time.Sleep(1 * time.Second)

	//for i := 0; i < 100; i++ {
	//	go generateUniqueSequence(i)
	//}
	//
	//time.Sleep(10 * time.Second)

}
