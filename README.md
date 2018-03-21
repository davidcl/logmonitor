# Simple HTTP log monitor

A simple console program that parse and monitor HTTP traffic from logs.

## Quick start

```
go get github.com/davidcl/logmonitor
go build github.com/davidcl/logmonitor
./logmonitor --help
```

## Dependencies

https://github.com/mattn/go-runewidth
https://github.com/nsf/termbox-go

## Test coverage

```
go test -coverprofile=coverage.out  github.com/davidcl/logmonitor
go tool cover -func=coverage.out
```

My focus while testing was to cover all "computation" behaviors while keeping the display / presentation logic untested.

## How to improve the design ?

The `maxHits` (statistics.go:133) could be reused accross calls to the `Display()` method. The hits counts are strictly ascending so the maxHits minimum could be stored to quickly filter out non-displayed values.

The `SortedInsert()` method put its values sorted into an array of pointers. Using a golang container/heap with its `Push()` method will slightly improve the performance. The drawback is that it will be harder to traverse the tree to retrieve the max hits counts.

`Counters` type (overruns.go:16) should be extended to compute content-length on average. This information is useful to evaluate the required bandwidth for the server.

Traffic counters (statistics.go:29) are stored into a doubly linked list to have a fixed time duration to compute overruns. 

## Note

This is my first Golang program; do not hesitate to point the Golang non-sense out !
