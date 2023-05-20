// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

// OriginalConfig sends the CONFIG GET ORIGINAL command to syslog-ng
func OriginalConfig(cc ControlChannel) (string, error) {
	return cc.SendCommand("CONFIG GET ORIGINAL")
}

// PreprocessedConfig sends the CONFIG GET PREPROCESSED command to syslog-ng
func PreprocessedConfig(cc ControlChannel) (string, error) {
	return cc.SendCommand("CONFIG GET PREPROCESSED")
}
