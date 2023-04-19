// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import (
	"errors"
	"io"
	"strings"
)

// Ping checks whether syslog-ng is listening on the control channel
func Ping(cc ControlChannel) error {
	rsp, err := cc.SendCommand("invalid") // we just send a command to check if syslog-ng is listening on the other end
	if strings.TrimSpace(rsp) == "" && errors.Is(err, io.EOF) {
		return nil
	}
	return err
}
