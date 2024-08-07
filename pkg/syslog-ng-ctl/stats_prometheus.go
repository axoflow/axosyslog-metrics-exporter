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
	"errors"
	"fmt"
	"strings"
	"time"

	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/prometheus/model/timestamp"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func createMetricsFromLegacyStats(legacyStats string) (map[string]*io_prometheus_client.MetricFamily, error) {
	var stats []Stat
	var err error

	stats, err = parseStats(legacyStats)
	if err != nil {
		return nil, err
	}

	mfs := make(map[string]*io_prometheus_client.MetricFamily)
	var errs []error
	const metric_ns = "syslogng"
	for _, stat := range stats {
		switch {
		case stat.SourceName == "global":
			switch stat.SourceID {
			case "scratch_buffers_count", "scratch_buffers_bytes":
				if err := pushMetric(mfs, metric_ns+"_"+stat.SourceID, io_prometheus_client.MetricType_GAUGE, nil, float64(stat.Number)); err != nil {
					errs = append(errs, err)
				}
			case "msg_allocated_bytes":
				if err := pushMetric(mfs, metric_ns+"_events_allocated_bytes", io_prometheus_client.MetricType_GAUGE, nil, float64(stat.Number)); err != nil {
					errs = append(errs, err)
				}
			default:
				// ignore other global stats
			}
		case stat.SourceName == "filter":
			labels := []*io_prometheus_client.LabelPair{
				newLabel("id", stat.SourceID),
				newLabel("result", stat.Type),
			}
			if err := pushMetric(mfs, metric_ns+"_filtered_events_total", io_prometheus_client.MetricType_COUNTER, labels, float64(stat.Number)); err != nil {
				errs = append(errs, err)
			}
		case strings.HasPrefix(stat.SourceName, "src.") && stat.SourceID != "" && stat.Type == "processed":
			labels := []*io_prometheus_client.LabelPair{
				newLabel("id", stat.SourceID),
				newLabel("result", stat.Type),
			}
			if stat.SourceInstance != "" {
				labels = append(labels, newLabel("driver_instance", stat.SourceInstance))
			}
			if err := pushMetric(mfs, metric_ns+"_input_events_total", io_prometheus_client.MetricType_COUNTER, labels, float64(stat.Number)); err != nil {
				errs = append(errs, err)
			}
		case strings.HasPrefix(stat.SourceName, "dst.") && slices.Contains([]string{"dropped", "queued", "written"}, stat.Type):
			result := stat.Type
			if result == "written" {
				result = "delivered"
			}
			labels := []*io_prometheus_client.LabelPair{
				newLabel("id", stat.SourceID),
				newLabel("result", result),
			}
			if stat.SourceInstance != "" {
				labels = append(labels, newLabel("driver_instance", stat.SourceInstance))
			}
			if err := pushMetric(mfs, metric_ns+"_output_events_total", io_prometheus_client.MetricType_COUNTER, labels, float64(stat.Number)); err != nil {
				errs = append(errs, err)
			}
		case stat.SourceName == "parser":
			labels := []*io_prometheus_client.LabelPair{
				newLabel("id", stat.SourceID),
				newLabel("result", stat.Type),
			}
			if err := pushMetric(mfs, metric_ns+"_parsed_events_total", io_prometheus_client.MetricType_COUNTER, labels, float64(stat.Number)); err != nil {
				errs = append(errs, err)
			}
		case stat.SourceName == "tag":
			labels := []*io_prometheus_client.LabelPair{
				newLabel("id", stat.SourceID),
				newLabel("result", stat.Type),
			}
			if err := pushMetric(mfs, metric_ns+"_tagged_events_total", io_prometheus_client.MetricType_COUNTER, labels, float64(stat.Number)); err != nil {
				errs = append(errs, err)
			}
		}
	}

	err = errors.Join(errs...)
	return mfs, err
}

func transformEventDelayMetric(delayMetric *io_prometheus_client.MetricFamily, delayMetricAge *io_prometheus_client.MetricFamily, now time.Time, lastMetricQueryTime time.Time, mfs map[string]*io_prometheus_client.MetricFamily) {

	if delayMetricAge == nil {
		delete(mfs, "syslogng_output_event_delay_sample_seconds")
		return
	}

	delayMetricAgeByLabel := make(map[string]*io_prometheus_client.Metric)
	for _, a := range delayMetricAge.Metric {
		delayMetricAgeByLabel[fmt.Sprint(a.Label)] = a
	}

	transformedMetric := []*io_prometheus_client.Metric{}
	for _, m := range delayMetric.Metric {
		delayMetric := m

		if d, ok := delayMetricAgeByLabel[fmt.Sprint(m.Label)]; ok {
			delayMetricAge := int(d.GetGauge().GetValue())

			lastDelaySampleTS := now.Add(time.Duration(-delayMetricAge) * time.Second)
			if lastDelaySampleTS.After(lastMetricQueryTime) {
				timestampMs := timestamp.FromTime(lastDelaySampleTS)
				transformedMetric = append(transformedMetric,
					&io_prometheus_client.Metric{
						Label:       delayMetric.GetLabel(),
						Gauge:       &io_prometheus_client.Gauge{Value: delayMetric.GetUntyped().Value},
						TimestampMs: &timestampMs,
					},
				)
			}
		}
	}

	if len(transformedMetric) == 0 {
		delete(mfs, "syslogng_output_event_delay_sample_seconds")
		return
	}

	delayMetric.Metric = transformedMetric
	delayMetric.Type = io_prometheus_client.MetricType_GAUGE.Enum()
}

// Workaround for a bug in older syslog-ng/AxoSyslog versions where the output of STATS PROMETHEUS was overescaped.
// Escapes \ as \\ everywhere except for the allowed sequences: \\, \n, \"
func sanitizeBuggyFormat(output string) string {
	var fixedOutput strings.Builder

	length := len(output)
	for i := 0; i < length; i++ {
		c := output[i]

		if c != '\\' {
			fixedOutput.WriteByte(c)
			continue
		}

		if i+1 >= length {
			fixedOutput.WriteString(`\\`)
			break
		}

		if next := output[i+1]; next == '\\' || next == 'n' || next == '"' {
			fixedOutput.WriteByte(c)
			fixedOutput.WriteByte(next)
			i++
			continue
		}

		fixedOutput.WriteString(`\\`)
	}

	return fixedOutput.String()
}

func StatsPrometheus(ctx context.Context, cc ControlChannel, lastMetricQueryTime *time.Time) ([]*io_prometheus_client.MetricFamily, error) {
	rsp, err := cc.SendCommand(ctx, "STATS PROMETHEUS")
	if err != nil {
		return nil, err
	}

	now := time.Now()
	defer func() { *lastMetricQueryTime = now }()

	var mfs map[string]*io_prometheus_client.MetricFamily
	if strings.HasPrefix(rsp, StatsHeader) {
		mfs, err = createMetricsFromLegacyStats(rsp)
		return maps.Values(mfs), err
	}

	rsp = sanitizeBuggyFormat(rsp)
	mfs, err = new(expfmt.TextParser).TextToMetricFamilies(strings.NewReader(rsp))

	var delayMetric *io_prometheus_client.MetricFamily
	var delayMetricAge *io_prometheus_client.MetricFamily

	for _, mf := range mfs {
		if mf.Type == nil {
			continue
		}
		if *mf.Type != io_prometheus_client.MetricType_UNTYPED {
			continue
		}
		if mf.Name == nil {
			continue
		}

		switch {
		case strings.HasSuffix(*mf.Name, "_events_total"):
			for _, m := range mf.Metric {
				m.Counter = &io_prometheus_client.Counter{
					Value: m.Untyped.Value,
				}
				m.Untyped = nil
			}
			mf.Type = io_prometheus_client.MetricType_COUNTER.Enum()
		case mf.GetName() == "syslogng_output_event_delay_sample_seconds":
			delayMetric = mf
		case mf.GetName() == "syslogng_output_event_delay_sample_age_seconds":
			delayMetricAge = mf
			fallthrough
		default:
			for _, m := range mf.Metric {
				m.Gauge = &io_prometheus_client.Gauge{
					Value: m.Untyped.Value,
				}
				m.Untyped = nil
			}
			mf.Type = io_prometheus_client.MetricType_GAUGE.Enum()
		}
	}

	if delayMetric != nil {
		transformEventDelayMetric(delayMetric, delayMetricAge, now, *lastMetricQueryTime, mfs)
	}

	return maps.Values(mfs), err
}

func pushMetric(mfs map[string]*io_prometheus_client.MetricFamily, name string, typ io_prometheus_client.MetricType, labels []*io_prometheus_client.LabelPair, value float64) error {
	m := &io_prometheus_client.Metric{
		Label: labels,
	}
	switch typ {
	case io_prometheus_client.MetricType_COUNTER:
		m.Counter = &io_prometheus_client.Counter{
			Value: &value,
		}
	case io_prometheus_client.MetricType_GAUGE:
		m.Gauge = &io_prometheus_client.Gauge{
			Value: &value,
		}
	default:
		return UnsupportedMetricType(typ)
	}

	mf := mfs[name]
	if mf == nil {
		mf = &io_prometheus_client.MetricFamily{
			Name: &name,
			Type: typ.Enum(),
		}
		mfs[name] = mf
	} else {
		if mf.Type == nil || *mf.Type != typ {
			return MetricTypeMismatch{
				ActualMetricFamily: mf,
				ExpectedType:       typ,
			}
		}
	}

	mf.Metric = append(mf.Metric, m)
	return nil
}

type MetricTypeMismatch struct {
	ActualMetricFamily *io_prometheus_client.MetricFamily
	ExpectedType       io_prometheus_client.MetricType
}

func (e MetricTypeMismatch) Error() string {
	name := "<nil>"
	if n := e.ActualMetricFamily.Name; n != nil {
		name = *n
	}
	return fmt.Sprintf("expected metric family %q to have type %q but had %q", name, e.ExpectedType, e.ActualMetricFamily.Type)
}

type UnsupportedMetricType io_prometheus_client.MetricType

func (e UnsupportedMetricType) Error() string {
	return fmt.Sprintf("metric type %q is not supported currently", io_prometheus_client.MetricType(e))
}

func newLabel(name string, value string) *io_prometheus_client.LabelPair {
	return &io_prometheus_client.LabelPair{
		Name:  &name,
		Value: &value,
	}
}
