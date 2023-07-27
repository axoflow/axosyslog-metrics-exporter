// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import (
	"context"
	"errors"
	"fmt"
	"strings"

	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func StatsPrometheus(ctx context.Context, cc ControlChannel) ([]*io_prometheus_client.MetricFamily, error) {
	rsp, err := cc.SendCommand(ctx, "STATS PROMETHEUS")
	if err != nil {
		return nil, err
	}

	var mfs map[string]*io_prometheus_client.MetricFamily
	if strings.HasPrefix(rsp, StatsHeader) {
		// received legacy-style stats
		var stats []Stat
		stats, err = parseStats(rsp)
		if err != nil {
			return nil, err
		}
		mfs = make(map[string]*io_prometheus_client.MetricFamily)
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
	} else {
		mfs, err = new(expfmt.TextParser).TextToMetricFamilies(strings.NewReader(rsp))
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
