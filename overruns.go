package main

import (
	"container/list"
	"time"
)

// OverrunEntry defines an out of threshold values
type OverrunEntry struct {
	hits  uint64
	start time.Time
	end   time.Time
}

// Counters describes all data used for statistics over time
type Counters struct {
	datetime time.Time // time of the counters
	hits     uint64
	//TODO add content-length
}

// Computes overruns and update traffic management
func UpdateOverruns(overruns []OverrunEntry, traffic *list.List, alertStart uint64, alertStop uint64) []OverrunEntry {

	// without traffic in the history, discard overrun detection and stop the in-progress overrun
	if traffic.Front() == nil {
		// stop the latest overrun !
		if len(overruns) > 0 {
			o := &overruns[len(overruns)-1]
			if o.start.Equal(o.end) {
				o.end = time.Now()
			}
		}
		return overruns
	}

	oldest := traffic.Front().Value.(Counters)
	newest := traffic.Back().Value.(Counters)

	// detect and update an in progress overrun
	if len(overruns) > 0 {
		o := &overruns[len(overruns)-1]
		if o.start.Equal(o.end) {
			// in progress overrun

			// stop if below the threshold
			if o.hits-oldest.hits <= alertStop {
				o.end = newest.datetime
				o.hits = newest.hits
			} else { // update otherwise
				o.hits = newest.hits
			}
			return overruns
		}
	}

	// no overrun is in progress, detect one
	wHits := newest.hits - oldest.hits
	if wHits >= alertStart {
		o := OverrunEntry{newest.hits, oldest.datetime, oldest.datetime}
		return append(overruns, o)
	}

	// nothing is detected
	return overruns
}

// Removes entries older than date (excluding)
func RemoveOlderThan(l *list.List, date time.Time) {
	var next *list.Element
	for oe := l.Front(); oe != nil; oe = next {
		t := oe.Value.(Counters).datetime
		if date.Equal(t) || date.Before(t) {
			break
		}

		next = oe.Next()
		l.Remove(oe)
	}
}
