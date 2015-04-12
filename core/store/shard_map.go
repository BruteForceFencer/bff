// Package store implements a map designed for concurrent use by splitting it
// into shards.
package store

import (
	"hash/fnv"
	"sync"
)

// NumShards is the number of shards for each ShardMap.
const NumShards = 128

// ShardMap is a concurrently available map.
type ShardMap struct {
	Type       string
	mutexes    []sync.Mutex
	shards     []map[interface{}]*BlockStatus
	maxTracked int64
}

// NewShardMap returns an initialized *ShardMap.
func NewShardMap(maxTracked int64) *ShardMap {
	result := new(ShardMap)
	result.maxTracked = maxTracked
	result.mutexes = make([]sync.Mutex, NumShards)
	result.shards = make([]map[interface{}]*BlockStatus, NumShards)
	for i := 0; i < NumShards; i++ {
		result.shards[i] = make(map[interface{}]*BlockStatus)
	}

	return result
}

// Get access the map and returns a *BlockStatus and a *Mutex.  If the mutex is
// not nil, you must unlock it when you're done.  If the value you pass to Get
// does not exist, it will be created for you.
func (s *ShardMap) Get(key interface{}) (*BlockStatus, *sync.Mutex) {
	index := s.hash(key)
	if index == -1 {
		return nil, nil
	}

	mutex := &s.mutexes[index]
	mutex.Lock()
	shard := s.shards[index]

	// Make sure we don't track too many values
	if s.maxTracked != 0 {
		totalApprox := NumShards * int64(len(shard))
		if totalApprox > s.maxTracked {
			mutex.Unlock()
			return nil, nil
		}
	}

	// Create a status if it doesn't exist
	status, exists := shard[key]
	if !exists {
		status = new(BlockStatus)
		shard[key] = status
	}

	return status, mutex
}

// hash returns the hash of the value such that it is between [0, NumShards).
func (s *ShardMap) hash(v interface{}) int {
	switch s.Type {
	case "string":
		val, ok := v.(string)
		if !ok {
			return -1
		}

		hash := fnv.New32()
		hash.Write([]byte(val))
		return int(hash.Sum32() % NumShards)
	case "int32":
		// Values that come as numbers in JSON are of type float64.
		val, ok := v.(float64)
		if !ok {
			return -1
		}

		return int(int32(val) % NumShards)
	default:
		return -1
	}
}

// CleanUp accesses every shard and deletes values whose FrontTile is less than
// the clock.
func (s *ShardMap) CleanUp(clock int32) {
	fClock := float64(clock)
	for i := range s.mutexes {
		m := &s.mutexes[i]
		m.Lock()

		for k := range s.shards[i] {
			if s.shards[i][k].FrontTile < fClock {
				delete(s.shards[i], k)
			}
		}

		m.Unlock()
	}
}

// Iterates calls f for every key-status pair in the ShardMap.
func (s *ShardMap) Iterate(f func(key interface{}, status *BlockStatus)) {
	for i := range s.mutexes {
		mutex := &s.mutexes[i]
		mutex.Lock()

		shard := s.shards[i]
		for key, status := range shard {
			f(key, status)
		}

		mutex.Unlock()
	}
}
