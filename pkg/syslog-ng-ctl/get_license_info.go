// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import "strings"

// GetLicenseInfo sends the LICENSE command to syslog-ng
func GetLicenseInfo(cc ControlChannel) (string, error) {
	info, err := cc.SendCommand("LICENSE")
	info = strings.TrimSpace(info)
	return info, err
}
