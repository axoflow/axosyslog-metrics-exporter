// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import "context"

// Reload sends the reload command to the syslog-ng instance behind the control channel
func Reload(ctx context.Context, cc ControlChannel) error {
	_, err := cc.SendCommand(ctx, "RELOAD")
	return err
}
