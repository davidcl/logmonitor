package main

import (
	"bytes"
	"fmt"
	"net"
	re "regexp"
	"strconv"
	"time"
)

// Common Log Format entry
type LogEntry struct {
	ip         net.IP
	ident      []byte
	userid     []byte
	datetime   time.Time
	request    []byte
	section    []byte
	resource   []byte
	protocol   []byte
	httpStatus int
	size       uint64
}

// Pretty print on errors
func (e LogEntry) String() string {
	return fmt.Sprintf("%s %s %s [%s] \"%s\" %d %d", e.ip, e.ident, e.userid, e.datetime.Format("02/Jan/2006:15:04:05 -0700"), e.request, e.httpStatus, e.size)

}

// Apache common log regexp
// TODO not RFC-1738 conformant, see "Reserved:" characters
var loggerRegexp = re.MustCompile("^(.+) (.+) (.+) \\[(.+)\\] \"(.+) /(.+/)?(.+)? (.+)\" ([0-9]+) ([0-9\\-]+).*$")

// Parse the log line into a specific variable
func Parse(line []byte) (*LogEntry, error) {
	found := loggerRegexp.FindSubmatchIndex(line)
	if found == nil {
		return nil, fmt.Errorf("invalid log entry : %s", line)
	}

	n := 2

	ip := net.ParseIP(string(line[found[n]:found[n+1]]))
	n += 2
	ident := line[found[n]:found[n+1]]
	n += 2
	userid := line[found[n]:found[n+1]]
	n += 2

	datetime, err := time.Parse("02/Jan/2006:15:04:05 -0700", string(line[found[n]:found[n+1]]))
	if err != nil {
		return nil, fmt.Errorf("invalid log DateTime format: %s", line[found[n]:found[n+1]])
	}
	n += 2

	request := line[found[n]:found[n+1]]
	n += 2
	section := []byte{}
	if found[n] > 0 {
		// only use the first "section"
		i := found[n] + bytes.IndexRune(line[found[n]:found[n+1]], '/')
		section = line[found[n]:i]
	}
	n += 2
	resource := []byte{}
	if found[n] > 0 {
		resource = line[found[n]:found[n+1]]
	}
	n += 2
	protocol := line[found[n]:found[n+1]]
	n += 2

	httpStatus, err := strconv.Atoi(string(line[found[n]:found[n+1]]))
	if err != nil {
		return nil, fmt.Errorf("invalid log HttpStatus number : %s", line[found[n]:found[n+1]])
	}
	n += 2

	size, err := strconv.ParseUint(string(line[found[n]:found[n+1]]), 10, 0)
	if err != nil {
		// on "-", the content length is set to 0
		size = 0
	}
	n += 2

	return &LogEntry{ip, ident, userid, datetime, request, section, resource, protocol, httpStatus, size}, nil
}
