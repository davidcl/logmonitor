package main

import (
	"testing"
)

func TestSingleLine(t *testing.T) {
	entry, err := Parse([]byte("127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] \"GET /my_section/apache_pb.gif HTTP/1.0\" 200 2326"))

	if err != nil {
		t.Fatal(err)
	}

	if entry.ip.String() != "127.0.0.1" {
		t.Fatal("invalid IP", entry)
	}
	if string(entry.ident) != "user-identifier" {
		t.Fatal("invalid ident", entry)
	}
	if string(entry.userid) != "frank" {
		t.Fatal("invalid userid", entry)
	}
	if entry.datetime.Unix() != 971211336 {
		t.Fatal("invalid datetime", entry)
	}
	if string(entry.request) != "GET" {
		t.Fatal("invalid request", entry)
	}
	if string(entry.section) != "my_section" {
		t.Fatal("invalid section", entry)
	}
	if string(entry.resource) != "apache_pb.gif" {
		t.Fatal("invalid resource", entry)
	}
	if string(entry.protocol) != "HTTP/1.0" {
		t.Fatal("invalid protocol", entry)
	}
	if entry.httpStatus != 200 {
		t.Fatal("invalid httpStatus", entry)
	}
	if entry.size != 2326 {
		t.Fatal("invalid size", entry)
	}
}

func TestWithoutSection(t *testing.T) {
	entry, err := Parse([]byte("127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326"))

	if err != nil {
		t.Fatal(err)
	}

	if len(entry.section) > 0 {
		t.Fatal("log without section not empty")
	}
}

func TestWithoutResource(t *testing.T) {
	entry, err := Parse([]byte("127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] \"GET / HTTP/1.0\" 200 2326"))

	if err != nil {
		t.Fatal(err)
	}

	if len(entry.section) > 0 {
		t.Fatal("log without section not empty")
	}
	if len(entry.resource) > 0 {
		t.Fatal("log without resource not empty")
	}
}

func TestWithMultipleSubdir(t *testing.T) {
	entry, err := Parse([]byte("127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] \"GET /foo/bar/baz/apache_pb.gif HTTP/1.0\" 200 2326"))

	if err != nil {
		t.Fatal(err)
	}

	if string(entry.section) != "foo" {
		t.Fatal("log with multiple subdir have wrong section")
	}
}
