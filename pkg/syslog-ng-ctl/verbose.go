// Copyright © 2023 Axoflow

package syslogngctl

// GetVerbose returns whether verbose logging is enabled on the syslog-ng instance behind the control channel
func GetVerbose(cc ControlChannel) (on bool, err error) {
	return getLog(cc, "VERBOSE")
}

// SetVerbose enables or disables verbose logging on the syslog-ng instance behind the control channel
func SetVerbose(cc ControlChannel, on bool) error {
	return setLog(cc, "VERBOSE", on)
}
