package osutil

import (
	"sync"
)

// SyncedMap holds a synchronized map.
type SyncedMap[K comparable, V any] struct {
	mutex sync.RWMutex
	m     map[K]V
}

// NewSyncedMap returns a nw synchronized map.
func NewSyncedMap[K comparable, V any]() (syncedMap *SyncedMap[K, V]) {
	return &SyncedMap[K, V]{
		m: map[K]V{},
	}
}

// Delete removes the entry with the given key from the map.
func (m *SyncedMap[K, V]) Delete(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.m, key)
}

// Get returns the value for the given key.
func (m *SyncedMap[K, V]) Get(key K) (value V, ok bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	value, ok = m.m[key]

	return value, ok
}

// Set sets the given value for the given key.
func (m *SyncedMap[K, V]) Set(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.m[key] = value
}
