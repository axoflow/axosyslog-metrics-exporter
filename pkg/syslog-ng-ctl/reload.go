// Copyright © 2023 Axoflow

package syslogngctl

// Reload sends the reload command to the syslog-ng instance behind the control channel
func Reload(cc ControlChannel) error {
	_, err := cc.SendCommand("RELOAD")
	return err
}
