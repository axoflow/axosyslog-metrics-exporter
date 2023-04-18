// Copyright © 2023 Axoflow
// All rights reserved.

package io

import (
	"io"
)

type nopRWCloser struct {
	io.ReadWriter
}

func (nopRWCloser) Close() error {
	return nil
}

func NopRWCloser(rw io.ReadWriter) io.ReadWriteCloser {
	return nopRWCloser{rw}
}
