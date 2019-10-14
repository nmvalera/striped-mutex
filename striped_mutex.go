package stripedmutex

import (
	"hash"
	"hash/fnv"
	"sync"
)

// StripedMutex is an object that allows fine grained locking based on keys
//
// It ensures that if `key1 == key2` then lock associated with `key1` is the same as the one associated with `key2`
// It holds a stable number of locks in memory that use can control
//
// It is inspired from Java lib Guava: https://github.com/google/guava/wiki/StripedExplained
type StripedMutex struct {
	stripes []*sync.Mutex
	pool    *sync.Pool
}

// Lock acquire lock for a given key
func (m *StripedMutex) Lock(key string) {
	l, _ := m.GetLock(key)
	l.Lock()
}

// Unlock release lock for a given key
func (m *StripedMutex) Unlock(key string) {
	l, _ := m.GetLock(key)
	l.Unlock()
}

// GetLock retrieve a lock for a given key
func (m *StripedMutex) GetLock(key string) (*sync.Mutex, error) {
	h := m.pool.Get().(hash.Hash64)
	defer m.pool.Put(h)
	h.Reset()
	_, err := h.Write([]byte(key))
	return m.stripes[h.Sum64()%uint64(len(m.stripes))], err
}

// New creates a StripedMutex
func New(stripes uint) *StripedMutex {
	m := &StripedMutex{
		make([]*sync.Mutex, stripes),
		&sync.Pool{New: func() interface{} { return fnv.New64() }},
	}
	for i := 0; i < len(m.stripes); i++ {
		m.stripes[i] = &sync.Mutex{}
	}

	return m
}
