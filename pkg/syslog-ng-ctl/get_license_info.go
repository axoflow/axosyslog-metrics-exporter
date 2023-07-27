// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import (
	"context"
	"strings"
)

// GetLicenseInfo sends the LICENSE command to syslog-ng
func GetLicenseInfo(ctx context.Context, cc ControlChannel) (string, error) {
	info, err := cc.SendCommand(ctx, "LICENSE")
	info = strings.TrimSpace(info)
	return info, err
}
