// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

// GetDebug returns whether debug logging is enabled on the syslog-ng instance behind the control channel
func GetDebug(cc ControlChannel) (on bool, err error) {
	return getLog(cc, "DEBUG")
}

// SetDebug enables or disables debug logging on the syslog-ng instance behind the control channel
func SetDebug(cc ControlChannel, on bool) error {
	return setLog(cc, "DEBUG", on)
}
