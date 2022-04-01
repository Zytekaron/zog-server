package cache

type Cache interface {
	// SetSize sets the size of the internal cache
	SetSize(size int)

	// Evict evicts a particular key from the cache
	Evict(id string)

	// Clear clears the internal cache
	Clear()
}
