// Package hitcounter augments the message-server with a store to track hits.
package hitcounter

import (
	"fmt"
	"github.com/BruteForceFencer/bff/logger"
	"time"
)

// HitCounter is a server that tracks several directions.
type HitCounter struct {
	Clock      *Clock
	Count      *RunningCount
	Directions map[string]*Direction
	Logger     *logger.Logger
}

// NewHitCounter returns an initialized *HitCounter.
func NewHitCounter(directions []Direction) *HitCounter {
	result := new(HitCounter)
	result.Clock = NewClock()
	result.Count = NewRunningCount(128, 24*time.Hour)

	result.Directions = make(map[string]*Direction)
	for i := range directions {
		result.Directions[directions[i].Name] = &directions[i]
		result.scheduleCleanUp(&directions[i])
	}

	return result
}

func (h *HitCounter) HandleRequest(direction string, value interface{}) bool {
	// Make sure the direction exists.
	dir, ok := h.Directions[direction]
	if !ok {
		return false
	}

	safe := dir.Hit(h.Clock.GetTime(), value)
	if !safe && h.Logger != nil {
		h.Logger.Log(direction, fmt.Sprint(value))
	}

	h.Count.Inc()
	return safe
}

func (h *HitCounter) scheduleCleanUp(dir *Direction) {
	go func(dir *Direction) {
		for {
			dir.Store.CleanUp(h.Clock.GetTime())
			time.Sleep(time.Duration(dir.CleanUpTime) * time.Second)
		}
	}(dir)
}
