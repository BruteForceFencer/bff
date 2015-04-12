// Package hitcounter augments the message-server with a store to track hits.
package hitcounter

import (
	"fmt"
	"github.com/BruteForceFencer/core/logger"
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
func NewHitCounter(directions []Direction, l *logger.Logger) *HitCounter {
	result := new(HitCounter)
	result.Clock = NewClock()
	result.Count = NewRunningCount(128, 24*time.Hour)
	result.Logger = l

	// We store the directions in a map instead of a slice for quick access.
	result.Directions = make(map[string]*Direction)
	for i := range directions {
		dir := &directions[i]

		result.Directions[dir.Name] = dir
		result.scheduleCleanUp(dir)
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
