package main

import (
	"container/list"
	"fmt"
	"testing"
	"time"
)

func TestRemoveOlderThan(t *testing.T) {
	traffic := list.New()

	t0, _ := time.Parse(time.Kitchen, "0:00AM")
	traffic.PushBack(Counters{t0, uint64(1)})

	delay, _ := time.ParseDuration("1m")
	t1 := t0.Add(delay)
	traffic.PushBack(Counters{t1, uint64(2)})

	RemoveOlderThan(traffic, t0)
	if traffic.Len() != 2 {
		t.Fatal("RemoveOlderThan t0")
	}

	RemoveOlderThan(traffic, t1)
	if traffic.Len() != 1 {
		t.Fatal("RemoveOlderThan t1")
	}

	// testing multiple values removal after a gap
	localTime := t1.Add(delay)

	localTime = localTime.Add(delay)
	traffic.PushBack(Counters{localTime, uint64(2)})

	localTime = localTime.Add(delay)
	traffic.PushBack(Counters{localTime, uint64(2)})

	RemoveOlderThan(traffic, localTime)
	if traffic.Len() != 1 {
		t.Fatal("RemoveOlderThan failed to remove multiple")
	}
}

func TestAppend(t *testing.T) {
	overruns := []OverrunEntry{}
	traffic := list.New()

	localTime := parseTime("0:00AM")
	hits := uint64(0)

	thresholdDelay, _ := time.ParseDuration("2m")
	delay, _ := time.ParseDuration("1m")
	largeDelay, _ := time.ParseDuration("3m")

	hits++
	traffic.PushBack(Counters{localTime, hits})
	RemoveOlderThan(traffic, localTime.Add(-thresholdDelay))
	overruns = UpdateOverruns(overruns, traffic, 2, 1)

	if len(overruns) != 0 {
		t.Fatal("overruns wrongly detected")
	}

	localTime = localTime.Add(delay)
	hits++
	traffic.PushBack(Counters{localTime, hits})
	RemoveOlderThan(traffic, localTime.Add(-thresholdDelay))
	overruns = UpdateOverruns(overruns, traffic, 2, 1)

	if len(overruns) != 0 {
		t.Fatal("overruns wrongly detected")
	}

	localTime = localTime.Add(delay)
	hits++
	traffic.PushBack(Counters{localTime, hits})
	RemoveOlderThan(traffic, localTime.Add(-thresholdDelay))
	overruns = UpdateOverruns(overruns, traffic, 2, 1)

	if len(overruns) != 1 {
		t.Fatal("overruns not detected")
	} else {
		oe := overruns[0]
		if oe.hits != 3 {
			t.Fatal("invalid overrun hits")
		}
		if oe.start != parseTime("0:00AM") {
			t.Fatal("invalid overrun start")
		}
		if oe.end != parseTime("0:00AM") {
			t.Fatal("invalid overrun stop")
		}
	}

	localTime = localTime.Add(thresholdDelay)
	hits++
	traffic.PushBack(Counters{localTime, hits})
	RemoveOlderThan(traffic, localTime.Add(-thresholdDelay))
	overruns = UpdateOverruns(overruns, traffic, 2, 1)

	if len(overruns) != 1 {
		t.Fatal("overruns not detected")
	} else {
		oe := overruns[0]
		if oe.hits != 4 {
			t.Fatal("invalid complete overrun hits")
		}
		if oe.start != parseTime("0:00AM") {
			t.Fatal("invalid complete overrun start")
		}
		if oe.end != parseTime("0:04AM") {
			t.Fatal("invalid complete overrun stop")
		}
	}

	localTime = localTime.Add(largeDelay)
	hits++
	traffic.PushBack(Counters{localTime, hits})
	RemoveOlderThan(traffic, localTime.Add(-thresholdDelay))
	overruns = UpdateOverruns(overruns, traffic, 2, 1)

	if len(overruns) != 1 {
		t.Fatal("wrongly detected overrun after large delay")
	}
}

func TestAppendWithoutTraffic(t *testing.T) {
	overruns := []OverrunEntry{
		OverrunEntry{1, parseTime("0:01AM"), parseTime("0:01AM")},
	}
	traffic := list.New()

	// the traffic is empty (cleared by RemoveOlderThan) but
	// an overrun is still in progress
	overruns = UpdateOverruns(overruns, traffic, 2, 1)

	if len(overruns) != 1 {
		t.Fatal("overruns not detected")
	} else {
		oe := overruns[0]
		if oe.hits != 1 {
			t.Fatal("invalid overrun hits")
		}
		if oe.start != parseTime("0:01AM") {
			t.Fatal("invalid overrun start")
		}
		// should be time.Now()
		if time.Now().Sub(oe.end).Seconds() > 1 {
			t.Fatal("invalid overrun stop")
		}
	}
}

func TestAppendAndRefreshHits(t *testing.T) {
	overruns := []OverrunEntry{
		OverrunEntry{1, parseTime("0:01AM"), parseTime("0:01AM")},
	}
	traffic := list.New()

	// an overrun is in progress and the traffic continue to raise
	hits := uint64(2)
	traffic.PushBack(Counters{parseTime("0:02AM"), hits})
	hits++
	traffic.PushBack(Counters{parseTime("0:03AM"), hits})
	hits++
	traffic.PushBack(Counters{parseTime("0:04AM"), hits})
	hits++
	traffic.PushBack(Counters{parseTime("0:05AM"), hits})

	overruns = UpdateOverruns(overruns, traffic, 2, 1)

	if len(overruns) != 1 {
		t.Fatal("overruns not detected")
	} else {
		oe := overruns[0]
		if oe.hits != hits {
			t.Fatal("invalid overrun hits")
		}
		if oe.start != parseTime("0:01AM") {
			t.Fatal("invalid overrun start")
		}
		if oe.end != parseTime("0:01AM") { // still in-progress after
			t.Fatal("invalid overrun stop")
		}
	}
}

// Generates a time.Time from a string Kitchen descriptor
func parseTime(str string) time.Time {
	t, _ := time.Parse(time.Kitchen, str)
	return t
}

// dumps both traffic and overruns, used to debug the test
func dump(traffic *list.List, overruns []OverrunEntry) {
	for e := traffic.Front(); e != nil; e = e.Next() {
		c := e.Value.(Counters)
		fmt.Println(c.hits, "-", c.datetime)
	}
	for i := 0; i < len(overruns); i++ {
		o := overruns[i]
		fmt.Println(o.hits, "-", o.start, o.end)
	}
}
