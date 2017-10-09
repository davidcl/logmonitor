package main

import (
	"sort"
)

// MaxHitEntry contains information for each displayed max hit
type MaxHitEntry struct {
	hits    uint64
	section string
}

// MaxHit is a sorted set with fixed size
// TODO investigate container/heap traversal, as heap interface (Less, Len, Swap) is implemented
type MaxHits []*MaxHitEntry

// Create a new max hits set
func MakeMaxHits(size int) MaxHits {
	maxEntries := make(MaxHits, size)
	for i := 0; i < size; i++ {
		maxEntries[i] = &MaxHitEntry{}
	}
	return maxEntries
}

func (m MaxHits) Len() int {
	return len(m)
}

func (m MaxHits) Less(i, j int) bool {
	// We want to pop the highest hits by using greater than
	return m[i].hits > m[j].hits
}

func (m MaxHits) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Insert a max hit entry keeping the order sorted
func (m MaxHits) SortedInsert(e *MaxHitEntry) {
	i := sort.Search(len(m), func(i int) bool { return m[i].hits <= e.hits })

	if i < len(m) {
		copy(m[i+1:], m[i:len(m)-1])
		m[i] = e
	}
}
