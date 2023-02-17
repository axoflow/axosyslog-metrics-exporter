package io

import (
	"bytes"
	"io"

	bytesx "github.com/axoflow/axo-edge/x/bytes"
)

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
			if cpl == len(b) { // b only grows the partial match
				continue
			}
			// sep could not be matched, revert
			bs = append(bs, sep[:matched]) // reintroduce the optimistically left out partial separator
			matched = 0
		}
		// no partial match
		if i := bytes.Index(b, sep); i != -1 { // b contains sep starting at i
			// split b after the separator
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

type ReadUntilOption func(*ReadUntilOptions)

type ReadUntilOptions struct {
	ReadBufferSize int
}
