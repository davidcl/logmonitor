# Simple HTTP log monitor

A simple console program that parse and monitor HTTP traffic from logs.

## Quick start

```
mkdir -p src/github.com/davidcl/logmonitor
tar -C src/github.com/davidcl/logmonitor -xzf logmonitor.tar.gz
export GOPATH=$(pwd):$GOPATH
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

Traffic counters (statistics.go:29) are stored into a doubly linked list to meet stated requirements (having a fixed time duration to compute overruns). If the requirements could be updated to "for the past 2 minutes or 100 log entries" then the counters could be implemented as a ring buffer (backed by an array) thus relaxing the allocation / deallocation / GC pressure on (probable) high server load.

## Note

This is my first Golang program, I really enjoyed the challenge of learning this new language in the limited timeframe ; do not hesitate to point the Golang non-sense out !
