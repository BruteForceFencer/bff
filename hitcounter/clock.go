package hitcounter

import (
	"sync/atomic"
	"time"
)

// Clock represents an integer clock that starts at 0.
type Clock struct {
	ticks  int32
	ticker <-chan time.Time
}

// NewClock returns an initialized and started *Clock.
func NewClock() *Clock {
	result := &Clock{
		ticker: time.Tick(time.Second),
	}

	// Start a goroutine to capture the ticks
	go func() {
		for {
			<-result.ticker
			atomic.AddInt32(&result.ticks, 1)
		}
	}()

	return result
}

// GetTime returns the number of seconds that have passed since intiialization.
func (c *Clock) GetTime() int32 {
	return atomic.LoadInt32(&c.ticks)
}
