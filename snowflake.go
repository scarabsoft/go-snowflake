package snowflake

import (
	"fmt"
	"github.com/scarabsoft/go-snowflake/internal"
)

type NodeIDProvider interface {
	internal.NodeIDProvider
}

func NewFixedNodeProvider(id uint8) NodeIDProvider {
	return internal.NewFixedNodeIdProvider(id)
}

// Generator generates a snowflake like ID which is unique if and only if the NodeIDProvider provides a unique ID
// It assumes that the provided clock makes progress, if the sequence exhausted the system will not continue producing IDs
// ID format:   |-----42 Epoch Bits-----|-----8 Node Bits-----|-----14 Sequence Bits-----|
type Generator interface {
	Next() (ID, error)
}

type generatorImpl struct {
	gen internal.SnowflakeGenerator
}

type idImpl struct {
	id uint64
}

func (i idImpl) ID() uint64 {
	return i.id
}

func (i idImpl) Weeks() uint64 {
	return i.Days() / 7
}

func (i idImpl) Days() uint64 {
	return i.Hours() / 24
}

func (i idImpl) Hours() uint64 {
	return i.Minutes() / 60
}

func (i idImpl) Minutes() uint64 {
	return i.Seconds() / 60
}

func (i idImpl) Seconds() uint64 {
	return i.ID() >> 22
}

func (i idImpl) NodeID() uint8 {
	return uint8(i.ID() >> (14))
}

func (i idImpl) Iteration() uint16 {
	return uint16(i.ID())
}

func (i idImpl) String() string {
	return fmt.Sprintf("%d", i.ID())
}

type ID interface {
	ID() uint64

	Weeks() uint64
	Days() uint64
	Hours() uint64
	Minutes() uint64
	Seconds() uint64

	NodeID() uint8
	Iteration() uint16
	String() string
}

func (g *generatorImpl) Next() (ID, error) {
	r, err := g.gen.Next()
	if err != nil {
		return nil, err
	}
	return &idImpl{r}, nil
}

// Clock provides a time in ms for the generator and will be called for every ID once
type Clock interface {
	internal.Clock
}

func NewUnixClock() Clock {
	return internal.NewUnixClockWithEpoch(0)
}

func NewUnixClockWithEpoch(epoch uint64) Clock {
	return internal.NewUnixClockWithEpoch(epoch)
}

type generatorBuilderImpl struct {
	clock        Clock
	nodeProvider NodeIDProvider
	maxSequence  uint16
}

type Option func(*generatorBuilderImpl) error

// WithClock sets a custom clock. Default system clock with UNIX epoch
func WithClock(clock Clock) Option {
	return func(impl *generatorBuilderImpl) error {
		impl.clock = clock
		return nil
	}
}

// WithNodeIDProvider sets the NodeIDProvider, which allows to generate nodeID based on hardware, like MAC or ...
// Make sure it generates a unique 8bit ID per node otherwise you will get duplicated IDs
func WithNodeIDProvider(provider NodeIDProvider) Option {
	return func(impl *generatorBuilderImpl) error {
		impl.nodeProvider = provider
		return nil
	}
}

// WithNodeID sets the id of the current Node. By default 1
func WithNodeID(nodeID uint8) Option {
	return func(impl *generatorBuilderImpl) error {
		impl.nodeProvider = NewFixedNodeProvider(nodeID)
		return nil
	}
}

// WithMaxSequence sets the max sequence per ms the system should support. By default 16,383 (16,383 ids can be generated per ms)
func WithMaxSequence(maxSeq uint16) Option {
	return func(impl *generatorBuilderImpl) error {
		impl.maxSequence = maxSeq
		return nil
	}
}

// New returns a new default generator and apply the requested options
//
// Default:
//		- Clock: system clock returning UNIX epoch
//		- Node: has ID 1
//		- MaxSequence: set to 16,383 (16,383 ids can be generated per ms)
func New(options ...Option) (Generator, error) {
	r := &generatorBuilderImpl{
		clock:        NewUnixClock(),
		nodeProvider: NewFixedNodeProvider(1),
		maxSequence:  internal.MaxSequence,
	}

	for _, option := range options {
		if err := option(r); err != nil {
			return nil, err
		}
	}

	seqProvider, err := internal.NewSequenceProvider(
		r.clock,
		r.maxSequence,
	)
	if err != nil {
		return nil, err
	}

	gen, err := internal.NewGenerator(
		seqProvider,
		r.nodeProvider,
	)

	if err != nil {
		return nil, err
	}

	return &generatorImpl{
		gen: gen,
	}, nil
}

func From(id uint64) ID {
	return &idImpl{id}
}
