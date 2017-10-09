package main

import (
	"fmt"
	"testing"
)

func TestSortedInsert(t *testing.T) {
	const N = 4
	maxHits := MakeMaxHits(N)

	for i := 1; i <= N; i++ {
		maxHits.SortedInsert(&MaxHitEntry{uint64(i), fmt.Sprintf("entry %d", i)})
	}
	for i := 1; i <= N; i++ {
		maxHits.SortedInsert(&MaxHitEntry{uint64(i), fmt.Sprintf("entry #2 %d", i)})
	}

	i := 0
	e := maxHits[i]
	if e.hits != 4 || e.section != "entry #2 4" {
		t.Fatal("invalid order 0")
	}
	i++
	e = maxHits[i]
	if e.hits != 4 || e.section != "entry 4" {
		t.Fatal("invalid order 1")
	}
	i++
	e = maxHits[i]
	if e.hits != 3 || e.section != "entry #2 3" {
		t.Fatal("invalid order 2")
	}
	i++
	e = maxHits[i]
	if e.hits != 3 || e.section != "entry 3" {
		t.Fatal("invalid order 3")
	}
}
