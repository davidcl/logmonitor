package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// Statistics provides a statistics container for all logged data.
type Statistics struct {
	sync.RWMutex // structure lock

	statisticsDuration time.Duration // const compute duration
	alertStart         uint64        // const start an overrun
	alertStop          uint64        // const end the overrun below this limit

	hitsCounters    map[string]uint64 // section hits
	latest          Counters          // current counters
	trafficCounters *list.List        // traffic history
	overruns        []OverrunEntry    // detected overruns (in-progress and stopped)
}

// Allocate a new Statistics instance
func NewStatistics(statisticsDuration time.Duration, alertStart uint64, alertStop uint64) *Statistics {
	s := new(Statistics)

	s.statisticsDuration = statisticsDuration
	s.alertStart = alertStart
	s.alertStop = alertStop

	s.hitsCounters = map[string]uint64{}
	s.trafficCounters = list.New()
	s.overruns = []OverrunEntry{}
	return s
}

// Process the argument, appending its content to the current statistics
func (s *Statistics) Process(r io.Reader) {
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			entry, err := Parse(scanner.Bytes())
			if err != nil {
				log.New(os.Stderr, "", 0).Fatal(err)
			}

			if entry != nil {
				s.Lock()
				// update hits
				hits := s.hitsCounters[string(entry.section)]
				hits++
				s.hitsCounters[string(entry.section)] = hits

				// update global counter
				s.latest.datetime = entry.datetime
				s.latest.hits++
				// store current counter and add an overrun entry when needed
				s.trafficCounters.PushBack(s.latest)
				// cleanup old values
				RemoveOlderThan(s.trafficCounters, entry.datetime)
				// update the overruns
				s.overruns = UpdateOverruns(s.overruns, s.trafficCounters, s.alertStart, s.alertStop)

				s.Unlock()

				// TODO store section if local hits count is more than the stored minimum of max hits
			}
		}
	}()
}

// Post-process and display to the console
func (s *Statistics) Display() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, height := termbox.Size()

	maxHits := MakeMaxHits(height)
	overruns := make([]OverrunEntry, height)

	var wg sync.WaitGroup
	wg.Add(2) // one per compute goroutines
	go s.compute_maxhits(&maxHits, &wg)
	go s.compute_overruns(&overruns, &wg)
	wg.Wait() // all data are present after this point

	// line to display shared index
	y := 0

	// overrun header
	for i := 0; i < len(overruns) && y < height; i, y = i+1, y+1 {
		o := overruns[i]
		if o.hits <= 0 {
			break
		}
		if o.end.After(o.start) {
			print_wide(0, y, fmt.Sprintf("High traffic generated an alert - hits = %d, triggered at %s, ended at %s", o.hits, o.start, o.end))
		} else {
			print_wide(0, y, fmt.Sprintf("High traffic generated an alert - hits = %d, triggered at %s", o.hits, o.start))
		}
	}

	// columns header
	if y < height {
		print_header(0, y, width, fmt.Sprintf("%s %s",
			"     HITS", "SECTION"))
		y++
	}

	// max hits
	for i := 0; i < len(maxHits) && maxHits[i].hits != 0 && y < height; i, y = i+1, y+1 {
		print_wide(0, y, fmt.Sprintf("%9.f %s", float64(maxHits[i].hits), maxHits[i].section))
	}
	termbox.Flush()
}

// compute the max hits per section
func (s Statistics) compute_maxhits(maxHits *MaxHits, wg *sync.WaitGroup) {
	defer wg.Done()
	s.RLock()
	// TODO the whole lookup and insert is greedy (but fast)
	// TODO share maxHits between calls
	for k, v := range s.hitsCounters {
		maxHits.SortedInsert(&MaxHitEntry{v, k})
	}
	s.RUnlock()
}

// compute the overruns
func (s Statistics) compute_overruns(overruns *[]OverrunEntry, wg *sync.WaitGroup) {
	defer wg.Done()

	s.RLock()
	// overruns are pre-computed, just copy them !
	copy(*overruns, s.overruns)
	s.RUnlock()
}

// print a string at a position
func print_wide(x, y int, s string) {
	for _, r := range s {
		termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)

		// from termbox : `CJK runes often require more than one cell to display them nicely`
		w := runewidth.RuneWidth(r)
		if w == 0 || (w == 2 && runewidth.IsAmbiguousWidth(r)) {
			w = 1
		}
		x += w
	}
}

// print the string with an header style
func print_header(x, y, width int, s string) {
	fg := termbox.ColorBlack
	bg := termbox.ColorWhite
	for _, r := range s {
		termbox.SetCell(x, y, r, fg, bg)

		// from termbox : `CJK runes often require more than one cell to display them nicely`
		w := runewidth.RuneWidth(r)
		if w == 0 || (w == 2 && runewidth.IsAmbiguousWidth(r)) {
			w = 1
		}
		x += w
	}
	for ; x < width; x++ {
		termbox.SetCell(x, y, ' ', fg, bg)
	}
}
