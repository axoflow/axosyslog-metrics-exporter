// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import (
	"context"
	"io"
	"net"
)

func NewUnixDomainSocketControlChannel(socketAddr string) ControlChannel {
	return NewReadWriterControlChannel(func(ctx context.Context) (io.ReadWriter, error) {
		var d net.Dialer
		return d.DialContext(ctx, "unix", socketAddr)
	})
}
