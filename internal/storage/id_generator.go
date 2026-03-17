package storage

// IDGenerator generates sequential, unique, monotonically increasing IDs.
type IDGenerator struct {
	counter int64
}

// NewIDGenerator returns a new IDGenerator with its counter initialised to 0.
func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

// NextID returns the next unique ID. Calls are thread-safe.
func (g *IDGenerator) NextID() int64 {
	// Stub — implementation to be added by the implementation developer.
	return 0
}
