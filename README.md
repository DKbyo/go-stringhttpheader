# go-stringhttpheader

go-stringhttpheader is a Go library for encoding structs into Header fields. Useful to create headers for [GCP Buckets Go library](https://pkg.go.dev/cloud.google.com/go/storage)

[![Build Status](https://github.com/dkbyo/go-stringhttpheader/workflows/CI/badge.svg?branch=master)](https://github.com/dkbyo/go-stringhttpheader/actions)
[![Coverage Status](https://coveralls.io/repos/github/dkbyo/go-stringhttpheader/badge.svg?branch=master)](https://coveralls.io/github/dkbyo/go-stringhttpheader?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/dkbyo/go-stringhttpheader)](https://goreportcard.com/report/github.com/dkbyo/go-stringhttpheader)
[![GoDoc](https://godoc.org/github.com/dkbyo/go-stringhttpheader?status.svg)](https://godoc.org/github.com/dkbyo/go-stringhttpheader)

## install

`go get -u github.com/dkbyo/go-stringhttpheader`


## usage

```go
package main

import (
	"fmt"
	"net/http"
	"github.com/dkbyo/go-stringhttpheader"
)

type Options struct {
	hide         string
	ContentType  string `header:"Content-Type"`
	Length       int
	XArray       []string `header:"X-Array"`
	TestHide     string   `header:"-"`
	IgnoreEmpty  string   `header:"X-Empty,omitempty"`
	IgnoreEmptyN string   `header:"X-Empty-N,omitempty"`
	CustomHeader http.Header
}

func main() {
	opt := Options{
		hide:         "hide",
		ContentType:  "application/json",
		Length:       2,
		XArray:       []string{"test1", "test2"},
		TestHide:     "hide",
		IgnoreEmptyN: "n",
		CustomHeader: http.Header{
			"X-Test-1": []string{"233"},
			"X-Test-2": []string{"666"},
		},
	}
	h, _ := stringhttpheader.Header(opt)
	fmt.Printf("%#v", h)
	// h:
	// string[]{
	//	"X-Test-1: 233",
	//	"X-Test-2: 666",
	//	"Content-Type: application/json",
	//	"Length: 2",
	//	"X-Array: test1, test2",
	//	"X-Empty-N: n",
	// }
	
}
```
