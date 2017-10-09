package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)
import (
	"github.com/nsf/termbox-go"
)

/*
* Options
 */
var delay time.Duration
var statisticsDuration time.Duration
var alertStart uint64
var alertStop uint64

// options declaration
func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... [LOG]...\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Simple program to monitor HTTP access LOGs files")
		fmt.Fprintln(os.Stderr, "")
		flag.PrintDefaults()
	}

	flag.DurationVar(&delay, "delay", time.Duration(10)*time.Second, "specifies the delay between screen updates")
	flag.DurationVar(&statisticsDuration, "statistics", time.Duration(2)*time.Minute, "specifies the duration used to compute statistics")
	flag.Uint64Var(&alertStart, "alert-start", 100, "generates an alert over this number of hits over the statistics duration")
	flag.Uint64Var(&alertStop, "alert-stop", 50, "ends the alert below this number of hits over the statistics duration")
}

func main() {
	flag.Parse()
	stats := NewStatistics(statisticsDuration, alertStart, alertStop)

	// process any file argument
	for _, arg := range flag.Args() {
		f, err := os.Open(arg)
		if err != nil {
			log.Fatalln(err)
		}
		stats.Process(bufio.NewReader(f))
	}

	// display the statistics
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	// manage redraw as termbox does not provide timers
	timer := time.Tick(delay)
	go func() {
		for {
			select {
			case <-timer:
				termbox.Interrupt()
			}
		}
	}()

	// main CLI loop
	stats.Display()
loop:
	for {

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			default:
				switch ev.Ch {
				case 'q':
					break loop
				}
			}
		case termbox.EventResize:
			stats.Display()
		case termbox.EventInterrupt: // timer
			stats.Display()
		default:
			break loop
		}

	}
}
