package hitcounter

import (
	"sync"
	"sync/atomic"
	"time"
)

// RunningCount is a counter where every increment is temporary.  After the
// duration has passed, the counter decrements.
type RunningCount struct {
	duration time.Duration

	// The lock is only used to access the index, not the count slice.  The
	// values of the count slice is only accessed using atomic operations.
	lock  sync.RWMutex
	index int
	count []uint64
}

// NewRunningCount returns an initialized *RunningCount.
func NewRunningCount(Granularity int, Duration time.Duration) *RunningCount {
	result := new(RunningCount)
	result.duration = Duration
	result.count = make([]uint64, Granularity)

	go func() {
		c := time.Tick(result.duration)
		for {
			<-c
			result.lock.Lock()
			result.index = (result.index + 1) % len(result.count)
			atomic.StoreUint64(&result.count[result.index], 0)
			result.lock.Unlock()
		}
	}()

	return result
}

// Inc increments the count by 1.  This increment will expire after the
// RunningCount's duration.
func (r *RunningCount) Inc() {
	r.lock.RLock()
	atomic.AddUint64(&r.count[r.index], 1)
	r.lock.RUnlock()
}

// Count returns the current count.
func (r *RunningCount) Count() uint64 {
	var result uint64

	for i := range r.count {
		result += atomic.LoadUint64(&r.count[i])
	}

	return result
}
