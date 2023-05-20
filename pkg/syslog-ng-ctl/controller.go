// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import io_prometheus_client "github.com/prometheus/client_model/go"

// Controller implements syslog-ng-ctl's functionality.
//
// Reference for available commands in syslog-ng-ctl's source code: https://github.com/syslog-ng/syslog-ng/blob/0e7c762c704efbda0ae10b61c35700ef0bdbb9c1/syslog-ng-ctl/syslog-ng-ctl.c#L111
type Controller struct {
	ControlChannel ControlChannel
}

func (c Controller) GetLicenseInfo() (string, error) {
	return GetLicenseInfo(c.ControlChannel)
}

func (c Controller) Ping() error {
	return Ping(c.ControlChannel)
}

func (c Controller) Reload() error {
	return Reload(c.ControlChannel)
}

func (c Controller) Stats() ([]Stat, error) {
	return Stats(c.ControlChannel)
}

func (c Controller) OriginalConfig() (string, error) {
	return OriginalConfig(c.ControlChannel)
}

func (c Controller) PreprocessedConfig() (string, error) {
	return PreprocessedConfig(c.ControlChannel)
}

func (c Controller) StatsPrometheus() ([]*io_prometheus_client.MetricFamily, error) {
	return StatsPrometheus(c.ControlChannel)
}
