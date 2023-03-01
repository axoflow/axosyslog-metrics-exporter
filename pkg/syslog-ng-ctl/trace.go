// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

// GetTrace returns whether trace logging is enabled on the syslog-ng instance behind the control channel
func GetTrace(cc ControlChannel) (on bool, err error) {
	return getLog(cc, "TRACE")
}

// SetTrace enables or disables trace logging on the syslog-ng instance behind the control channel
func SetTrace(cc ControlChannel, on bool) error {
	return setLog(cc, "TRACE", on)
}
