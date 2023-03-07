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

func (c Controller) GetDebug() (on bool, err error) {
	return GetDebug(c.ControlChannel)
}

func (c Controller) GetTrace() (on bool, err error) {
	return GetTrace(c.ControlChannel)
}

func (c Controller) GetVerbose() (on bool, err error) {
	return GetVerbose(c.ControlChannel)
}

func (c Controller) Reload() error {
	return Reload(c.ControlChannel)
}

func (c Controller) SetDebug(on bool) error {
	return SetDebug(c.ControlChannel, on)
}

func (c Controller) SetTrace(on bool) error {
	return SetTrace(c.ControlChannel, on)
}

func (c Controller) SetVerbose(on bool) error {
	return SetVerbose(c.ControlChannel, on)
}

func (c Controller) Stats() ([]Stat, error) {
	return Stats(c.ControlChannel)
}

func (c Controller) StatsPrometheus() ([]io_prometheus_client.MetricFamily, error) {
	return StatsPrometheus(c.ControlChannel)
}
