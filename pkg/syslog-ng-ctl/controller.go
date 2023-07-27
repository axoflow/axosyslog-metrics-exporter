// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import (
	"context"

	io_prometheus_client "github.com/prometheus/client_model/go"
)

// Controller implements syslog-ng-ctl's functionality.
//
// Reference for available commands in syslog-ng-ctl's source code: https://github.com/syslog-ng/syslog-ng/blob/0e7c762c704efbda0ae10b61c35700ef0bdbb9c1/syslog-ng-ctl/syslog-ng-ctl.c#L111
type Controller struct {
	ControlChannel ControlChannel
}

func (c Controller) GetLicenseInfo(ctx context.Context) (string, error) {
	return GetLicenseInfo(ctx, c.ControlChannel)
}

func (c Controller) Ping(ctx context.Context) error {
	return Ping(ctx, c.ControlChannel)
}

func (c Controller) Reload(ctx context.Context) error {
	return Reload(ctx, c.ControlChannel)
}

func (c Controller) Stats(ctx context.Context) ([]Stat, error) {
	return Stats(ctx, c.ControlChannel)
}

func (c Controller) OriginalConfig(ctx context.Context) (string, error) {
	return OriginalConfig(ctx, c.ControlChannel)
}

func (c Controller) PreprocessedConfig(ctx context.Context) (string, error) {
	return PreprocessedConfig(ctx, c.ControlChannel)
}

func (c Controller) StatsPrometheus(ctx context.Context) ([]*io_prometheus_client.MetricFamily, error) {
	return StatsPrometheus(ctx, c.ControlChannel)
}
