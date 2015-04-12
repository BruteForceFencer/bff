package hitcounter

import (
	"github.com/BruteForceFencer/core/store"
)

// Direction is a tracked resource.  See package config for definition of its
// member variables.
type Direction struct {
	Store       *store.ShardMap
	Name        string
	CleanUpTime float64
	MaxHits     float64
	WindowSize  float64
}

// Hit registers a used value.
func (d *Direction) Hit(clock int32, val interface{}) bool {
	status, mutex := d.Store.Get(val)
	if status == nil {
		return false
	}
	defer mutex.Unlock()

	fClock := float64(clock)
	threshold := fClock + d.WindowSize

	if status.FrontTile < fClock {
		// The window was clear.
		status.FrontTile = fClock
		status.IsBlocked = false
	}

	status.FrontTile += d.WindowSize / d.MaxHits

	if status.FrontTile >= threshold {
		// We're above the threshold, but this happened only after incrementing
		// the front tile.

		if !status.IsBlocked {
			// This is the start time for this block.  We need to record it.
			status.Since = fClock
		}

		status.FrontTile = threshold
		status.IsBlocked = true
		return false
	} else {
		return !status.IsBlocked
	}
}

// Returns the list of all values in the map that have IsBlocked == true.
func (d *Direction) BlockedValues() []BlockedValue {
	result := make([]BlockedValue, 0)

	d.Store.Iterate(func(key interface{}, status *store.BlockStatus) {
		if status.IsBlocked {
			result = append(
				result,
				BlockedValue{Since: status.Since, Value: key},
			)
		}
	})

	return result
}
