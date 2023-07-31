// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import (
	"context"
	"errors"
	"io"
	"strings"
)

// Ping checks whether syslog-ng is listening on the control channel
func Ping(ctx context.Context, cc ControlChannel) error {
	rsp, err := cc.SendCommand(ctx, "invalid") // we just send a command to check if syslog-ng is listening on the other end
	if strings.TrimSpace(rsp) == "" && errors.Is(err, io.EOF) {
		return nil
	}
	return err
}
