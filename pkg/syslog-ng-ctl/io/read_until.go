// Copyright © 2023 Axoflow
// All rights reserved.

package io

import (
	"bytes"
	"io"

	bytesx "github.com/axoflow/metrics-exporter/pkg/syslog-ng-ctl/bytes"
)

// ReadUntil reads from the specified reader until it reaches the specified separator (or an error occurs, which includes EOF).
// It returns all bytes read until the separator, the rest of the bytes which were already read, and any error that happened.
func ReadUntil(rdr io.Reader, sep []byte, opts ...ReadUntilOption) (res []byte, rst []byte, err error) {
	options := ReadUntilOptions{
		ReadBufferSize: 4096,
	}
	for _, opt := range opts {
		opt(&options)
	}

	var bs [][]byte
	matched := 0
loop:
	for {
		b := make([]byte, options.ReadBufferSize)
		n, e := rdr.Read(b)
		b = b[:n] // limit b to its effective portion

		if matched > 0 { // we already have a partial match
			cpl := bytesx.CommonPrefixLen(b, sep[matched:])
			matched += cpl
			if matched == len(sep) { // sep is fully matched
				rst = append(rst, sep...)
				rst = append(rst, b[cpl:]...)
				break
			}
			if cpl > 0 && cpl == len(b) { // b only grows the partial match
				continue
			}
			// sep could not be matched, revert
			bs = append(bs, sep[:matched-cpl]) // reintroduce the optimistically left out partial separator
			matched = 0
		}
		// no partial match
		if i := bytes.Index(b, sep); i != -1 { // b contains sep starting at i
			// split b before the separator
			bs = append(bs, b[:i])
			rst = b[i:]
			break
		}
		// b does not contain the whole separator
		// try matching increasingly shorter prefixes of sep at the end of b
		for l := len(sep) - 1; l > 0; l-- {
			if bytes.HasSuffix(b, sep[:l]) { // partial match at the end of the current buffer
				matched = l
				bs = append(bs, b[:n-l])
				continue loop
			}
		}
		// b does not contain sep
		bs = append(bs, b)
		if e != nil {
			err = e
			break
		}
	}
	res = bytes.Join(bs, nil)
	return
}

// ReadUntilOption is an option for ReadUntil
type ReadUntilOption func(*ReadUntilOptions)

// ReadUntilOptions are the available options for ReadUntil
type ReadUntilOptions struct {
	// ReadBufferSize is the size of the buffer passed to the reader (default: 4096)
	ReadBufferSize int
}
