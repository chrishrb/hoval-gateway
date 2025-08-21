package inmemory

import (
	"sync"

	"github.com/chrishrb/hoval-gateway/hoval"
	"k8s.io/utils/clock"
)

// Store is an in-memory implementation of the store.Engine interface. As everything
// is stored in memory it is not stateless and cannot be used if running >1 instances.
// It is primarily provided to support unit testing.
type Store struct {
	sync.Mutex
	clock      clock.PassiveClock
	datapoints map[string]*hoval.Datapoint
	devices    map[uint32]*hoval.Device
}

func NewStore(clock clock.PassiveClock) *Store {
	return &Store{
		clock:      clock,
		datapoints: make(map[string]*hoval.Datapoint),
		devices:    make(map[uint32]*hoval.Device),
	}
}
