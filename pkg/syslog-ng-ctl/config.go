// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import "context"

// OriginalConfig sends the CONFIG GET ORIGINAL command to syslog-ng
func OriginalConfig(ctx context.Context, cc ControlChannel) (string, error) {
	return cc.SendCommand(ctx, "CONFIG GET ORIGINAL")
}

// PreprocessedConfig sends the CONFIG GET PREPROCESSED command to syslog-ng
func PreprocessedConfig(ctx context.Context, cc ControlChannel) (string, error) {
	return cc.SendCommand(ctx, "CONFIG GET PREPROCESSED")
}
