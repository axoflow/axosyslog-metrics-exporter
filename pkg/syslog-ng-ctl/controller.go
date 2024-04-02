// Copyright © 2023 Axoflow
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package syslogngctl

import (
	"context"
	"time"

	io_prometheus_client "github.com/prometheus/client_model/go"
)

// Controller implements syslog-ng-ctl's functionality.
//
// Reference for available commands in syslog-ng-ctl's source code: https://github.com/syslog-ng/syslog-ng/blob/0e7c762c704efbda0ae10b61c35700ef0bdbb9c1/syslog-ng-ctl/syslog-ng-ctl.c#L111
type Controller struct {
	ControlChannel      ControlChannel
	lastMetricQueryTime time.Time
}

func NewController(controlChannel ControlChannel) *Controller {
	return &Controller{
		ControlChannel:      controlChannel,
		lastMetricQueryTime: time.Now(),
	}
}

func (c *Controller) GetLicenseInfo(ctx context.Context) (string, error) {
	return GetLicenseInfo(ctx, c.ControlChannel)
}

func (c *Controller) Ping(ctx context.Context) error {
	return Ping(ctx, c.ControlChannel)
}

func (c *Controller) Reload(ctx context.Context) error {
	return Reload(ctx, c.ControlChannel)
}

func (c *Controller) Stats(ctx context.Context) ([]Stat, error) {
	return Stats(ctx, c.ControlChannel)
}

func (c *Controller) OriginalConfig(ctx context.Context) (string, error) {
	return OriginalConfig(ctx, c.ControlChannel)
}

func (c *Controller) PreprocessedConfig(ctx context.Context) (string, error) {
	return PreprocessedConfig(ctx, c.ControlChannel)
}

func (c *Controller) StatsPrometheus(ctx context.Context) ([]*io_prometheus_client.MetricFamily, error) {
	return StatsPrometheus(ctx, c.ControlChannel, &c.lastMetricQueryTime)
}
