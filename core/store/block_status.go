package store

// BlockStatus represents the status of a tracked value in a ShardMap.
type BlockStatus struct {
	IsBlocked bool

	// Frontile is a number that represents how frequently the value has been
	// hit.  Given a max number of hits x within a window of y seconds, each
	// hit will increment the FrontTile by y / x.
	FrontTile float64

	// Since is the number of seconds that have elapsed since the value was
	// blocked.
	Since float64
}
