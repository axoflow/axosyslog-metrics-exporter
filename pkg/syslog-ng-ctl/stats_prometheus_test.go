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
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

func TestStatsPrometheus(t *testing.T) {
	expected := []*io_prometheus_client.MetricFamily{
		{
			Name: amp("syslogng_events_allocated_bytes"),
			Type: io_prometheus_client.MetricType_GAUGE.Enum(),
			Metric: []*io_prometheus_client.Metric{
				{
					Gauge: &io_prometheus_client.Gauge{
						Value: amp(0.0),
					},
				},
			},
		},
		{
			Name: amp("syslogng_scratch_buffers_bytes"),
			Type: io_prometheus_client.MetricType_GAUGE.Enum(),
			Metric: []*io_prometheus_client.Metric{
				{
					Gauge: &io_prometheus_client.Gauge{
						Value: amp(0.0),
					},
				},
			},
		},
		{
			Name: amp("syslogng_scratch_buffers_count"),
			Type: io_prometheus_client.MetricType_GAUGE.Enum(),
			Metric: []*io_prometheus_client.Metric{
				{
					Gauge: &io_prometheus_client.Gauge{
						Value: amp(2.0),
					},
				},
			},
		},
		{
			Name: amp("syslogng_filtered_events_total"),
			Type: io_prometheus_client.MetricType_COUNTER.Enum(),
			Metric: []*io_prometheus_client.Metric{
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", "ff"),
						newLabel("result", "matched"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", "ff"),
						newLabel("result", "not_matched"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", "#anon-filter0"),
						newLabel("result", "matched"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", "#anon-filter0"),
						newLabel("result", "not_matched"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
			},
		},
		{
			Name: amp("syslogng_input_events_total"),
			Type: io_prometheus_client.MetricType_COUNTER.Enum(),
			Metric: []*io_prometheus_client.Metric{
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", "#anon-source0#0"),
						newLabel("driver_instance", "-"),
						newLabel("result", "processed"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", "s_network#1"),
						newLabel("result", "processed"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
			},
		},
		{
			Name: amp("syslogng_output_events_total"),
			Type: io_prometheus_client.MetricType_COUNTER.Enum(),
			Metric: []*io_prometheus_client.Metric{
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", "d_dest#0"),
						newLabel("driver_instance", "tcp,127.0.0.1:5555"),
						newLabel("result", "dropped"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", "d_dest#0"),
						newLabel("driver_instance", "tcp,127.0.0.1:5555"),
						newLabel("result", "queued"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", "d_dest#0"),
						newLabel("driver_instance", "tcp,127.0.0.1:5555"),
						newLabel("result", "delivered"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", "d_dest#1"),
						newLabel("driver_instance", "http,https://localhost:8080"),
						newLabel("result", "dropped"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", "d_dest#1"),
						newLabel("driver_instance", "http,https://localhost:8080"),
						newLabel("result", "queued"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", "d_dest#1"),
						newLabel("driver_instance", "http,https://localhost:8080"),
						newLabel("result", "delivered"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
			},
		},
		{
			Name: amp("syslogng_parsed_events_total"),
			Type: io_prometheus_client.MetricType_COUNTER.Enum(),
			Metric: []*io_prometheus_client.Metric{
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", "#anon-parser0"),
						newLabel("result", "processed"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", "#anon-parser0"),
						newLabel("result", "discarded"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
			},
		},
		{
			Name: amp("syslogng_tagged_events_total"),
			Type: io_prometheus_client.MetricType_COUNTER.Enum(),
			Metric: []*io_prometheus_client.Metric{
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", ".source.s_network"),
						newLabel("result", "processed"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("id", ".source.#anon-source0"),
						newLabel("result", "processed"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(0.0),
					},
				},
			},
		},
	}
	sortMetricFamilies(expected)

	expectedDelayMetrics := []*io_prometheus_client.MetricFamily{
		{
			Name: amp("syslogng_output_event_delay_sample_seconds"),
			Type: io_prometheus_client.MetricType_GAUGE.Enum(),
			Metric: []*io_prometheus_client.Metric{
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("driver", "http"),
						newLabel("url", "http://localhost/asd"),
						newLabel("id", "#anon-destination0#1"),
						newLabel("worker", "0"),
					},
					Gauge: &io_prometheus_client.Gauge{
						Value: amp(2.0),
					},
				},
			},
		},
		{
			Name: amp("syslogng_output_event_delay_sample_age_seconds"),
			Type: io_prometheus_client.MetricType_GAUGE.Enum(),
			Metric: []*io_prometheus_client.Metric{
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("driver", "http"),
						newLabel("url", "http://localhost/asd"),
						newLabel("id", "#anon-destination0#1"),
						newLabel("worker", "0"),
					},
					Gauge: &io_prometheus_client.Gauge{
						Value: amp(1.0),
					},
				},
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("transport", "tcp"),
						newLabel("address", "localhost:5555"),
						newLabel("driver", "afsocket"),
						newLabel("id", "#anon-destination0#0"),
					},
					Gauge: &io_prometheus_client.Gauge{
						Value: amp(31.0),
					},
				},
			},
		},
	}
	sortMetricFamilies(expectedDelayMetrics)

	expectedEscapeMetrics := []*io_prometheus_client.MetricFamily{
		{
			Name: amp("syslogng_classified_output_events_total"),
			Type: io_prometheus_client.MetricType_COUNTER.Enum(),
			Metric: []*io_prometheus_client.Metric{
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("app", "MSWinEventLog\\t1\\tSecurity\\t921448325\\tFri"),
						newLabel("source", "s_critical_hosts_515"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(1.0),
					},
				},
				{
					Label: []*io_prometheus_client.LabelPair{
						newLabel("app", "\\a\\t\n\"\\xfa\\"),
						newLabel("source", "s_unescaped_bug"),
					},
					Counter: &io_prometheus_client.Counter{
						Value: amp(1.0),
					},
				},
			},
		},
	}
	sortMetricFamilies(expectedEscapeMetrics)

	testCases := map[string]struct {
		cc       ControlChannel
		expected []*io_prometheus_client.MetricFamily
	}{
		"syslog-ng-ctl stats response for stats prometheus request": {
			cc: ControlChannelFunc(func(_ context.Context, cmd string) (rsp string, err error) {
				require.Equal(t, "STATS PROMETHEUS", cmd)
				return LEGACY_STATS_OUTPUT, nil
			}),
			expected: expected,
		},
		"syslog-ng stats prometheus response": {
			cc: ControlChannelFunc(func(_ context.Context, cmd string) (rsp string, err error) {
				require.Equal(t, "STATS PROMETHEUS", cmd)
				return PROMETHEUS_METRICS_OUTPUT, nil
			}),
			expected: expected,
		},
		"syslog-ng stats prometheus delay metrics": {
			cc: ControlChannelFunc(func(_ context.Context, cmd string) (rsp string, err error) {
				require.Equal(t, "STATS PROMETHEUS", cmd)
				return PROMETHEUS_DELAY_METRICS_OUTPUT, nil
			}),
			expected: expectedDelayMetrics,
		},
		"syslog-ng stats prometheus label escaping": {
			cc: ControlChannelFunc(func(_ context.Context, cmd string) (rsp string, err error) {
				require.Equal(t, "STATS PROMETHEUS", cmd)
				return PROMETHEUS_ESCAPE_METRICS_OUTPUT, nil
			}),
			expected: expectedEscapeMetrics,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			lastMetricQueryTime := time.Now().Add(-time.Second * 30)
			res, err := StatsPrometheus(context.Background(), testCase.cc, &lastMetricQueryTime)
			require.NoError(t, err)
			sortMetricFamilies(res)
			removeTimestamps(res)
			if !cmp.Equal(t, testCase.expected, cmpopts.IgnoreUnexported()) {
				assert.Equal(t, metricFamiliesToText(testCase.expected), metricFamiliesToText(res))
			}
		})
	}
}

type ControlChannelFunc func(ctx context.Context, cmd string) (rsp string, err error)

func (fn ControlChannelFunc) SendCommand(ctx context.Context, cmd string) (rsp string, err error) {
	return fn(ctx, cmd)
}

const LEGACY_STATS_OUTPUT = `SourceName;SourceId;SourceInstance;State;Type;Number
global;scratch_buffers_count;;a;queued;2
src.facility;;18;a;processed;0
src.stdin;#anon-source0#0;-;a;msg_size_max;0
src.facility;;7;a;processed;0
dst.network;d_dest#0;tcp,127.0.0.1:5555;a;eps_since_start;0
dst.network;d_dest#0;tcp,127.0.0.1:5555;a;msg_size_max;0
src.stdin;#anon-source0#0;-;a;processed;0
src.stdin;#anon-source0#0;-;a;stamp;0
dst.network;d_dest#0;tcp,127.0.0.1:5555;a;eps_last_1h;0
source;#anon-source0;;a;processed;0
dst.network;d_dest#0;tcp,127.0.0.1:5555;a;processed;0
src.facility;;19;a;processed;0
global;sdata_updates;;a;processed;0
src.facility;;8;a;processed;0
dst.network;d_dest#0;tcp,127.0.0.1:5555;a;msg_size_avg;0
dst.http;d_dest#1;http,https://localhost:8080;a;batch_size_max;0
center;;received;a;processed;0
dst.network;d_dest#0;tcp,127.0.0.1:5555;a;memory_usage;0
src.facility;;20;a;processed;0
src.facility;;9;a;processed;0
dst.network;d_dest#0;tcp,127.0.0.1:5555;a;truncated_bytes;0
center;;queued;a;processed;0
dst.http;d_dest#1;http,https://localhost:8080;a;batch_size_avg;0
src.facility;;other;a;processed;0
src.facility;;21;a;processed;0
dst.network;d_dest#0;tcp,127.0.0.1:5555;a;truncated_count;0
tag;.source.s_network;;a;processed;0
src.facility;;10;a;processed;0
dst.network;d_dest#0;tcp,127.0.0.1:5555;a;dropped;0
dst.network;d_dest#0;tcp,127.0.0.1:5555;a;queued;0
dst.network;d_dest#0;tcp,127.0.0.1:5555;a;written;0
src.severity;;0;a;processed;0
src.facility;;22;a;processed;0
src.facility;;11;a;processed;0
src.facility;;0;a;processed;0
src.severity;;1;a;processed;0
src.facility;;23;a;processed;0
filter;ff;;a;matched;0
filter;ff;;a;not_matched;0
src.facility;;12;a;processed;0
src.network;s_network;afsocket_sd.(stream,AF_INET(0.0.0.0:4444));a;connections;0
global;msg_clones;;a;processed;0
src.facility;;1;a;processed;0
src.severity;;2;a;processed;0
destination;d_dest;;a;processed;0
global;msg_allocated_bytes;;a;value;0
dst.network;d_dest#0;tcp,127.0.0.1:5555;a;eps_last_24h;0
src.facility;;13;a;processed;0
src.severity;;3;a;processed;0
src.facility;;2;a;processed;0
dst.http;d_dest#1;http,https://localhost:8080;a;eps_last_24h;0
dst.http;d_dest#1;http,https://localhost:8080;a;memory_usage;0
dst.http;d_dest#1;http,https://localhost:8080;a;dropped;0
dst.http;d_dest#1;http,https://localhost:8080;a;queued;0
dst.http;d_dest#1;http,https://localhost:8080;a;written;0
src.facility;;14;a;processed;0
dst.http;d_dest#1;http,https://localhost:8080;a;msg_size_max;0
src.facility;;3;a;processed;0
source;s_network;;a;processed;0
dst.http;d_dest#1;http,https://localhost:8080;a;eps_last_1h;0
tag;.source.#anon-source0;;a;processed;0
src.severity;;4;a;processed;0
dst.http;d_dest#1;http,https://localhost:8080;a;processed;0
src.facility;;15;a;processed;0
src.facility;;4;a;processed;0
src.severity;;5;a;processed;0
dst.http;d_dest#1;http,https://localhost:8080;a;msg_size_avg;0
global;payload_reallocs;;a;processed;0
global;scratch_buffers_bytes;;a;queued;0
src.facility;;16;a;processed;0
parser;#anon-parser0;;a;processed;0
parser;#anon-parser0;;a;discarded;0
filter;#anon-filter0;;a;matched;0
filter;#anon-filter0;;a;not_matched;0
src.facility;;5;a;processed;0
src.severity;;6;a;processed;0
src.stdin;#anon-source0#0;-;a;eps_last_1h;0
src.stdin;#anon-source0#0;-;a;eps_since_start;0
src.internal;s_network#1;;a;processed;0
src.internal;s_network#1;;a;stamp;0
src.stdin;#anon-source0#0;-;a;msg_size_avg;0
src.severity;;7;a;processed;0
src.facility;;6;a;processed;0
src.facility;;17;a;processed;0
dst.http;d_dest#1;http,https://localhost:8080;a;eps_since_start;0
src.stdin;#anon-source0#0;-;a;eps_last_24h;0
`
const PROMETHEUS_METRICS_OUTPUT = `syslogng_events_allocated_bytes 0
syslogng_filtered_events_total{id="#anon-filter0",result="matched"} 0
syslogng_filtered_events_total{id="#anon-filter0",result="not_matched"} 0
syslogng_filtered_events_total{id="ff",result="matched"} 0
syslogng_filtered_events_total{id="ff",result="not_matched"} 0
syslogng_input_events_total{id="#anon-source0#0",driver_instance="-",result="processed"} 0
syslogng_input_events_total{id="s_network#1",result="processed"} 0
syslogng_output_events_total{id="d_dest#0",driver_instance="tcp,127.0.0.1:5555",result="delivered"} 0
syslogng_output_events_total{id="d_dest#0",driver_instance="tcp,127.0.0.1:5555",result="dropped"} 0
syslogng_output_events_total{id="d_dest#0",driver_instance="tcp,127.0.0.1:5555",result="queued"} 0
syslogng_output_events_total{id="d_dest#1",driver_instance="http,https://localhost:8080",result="delivered"} 0
syslogng_output_events_total{id="d_dest#1",driver_instance="http,https://localhost:8080",result="dropped"} 0
syslogng_output_events_total{id="d_dest#1",driver_instance="http,https://localhost:8080",result="queued"} 0
syslogng_parsed_events_total{id="#anon-parser0",result="discarded"} 0
syslogng_parsed_events_total{id="#anon-parser0",result="processed"} 0
syslogng_scratch_buffers_bytes 0
syslogng_scratch_buffers_count 2
syslogng_tagged_events_total{id=".source.#anon-source0",result="processed"} 0
syslogng_tagged_events_total{id=".source.s_network",result="processed"} 0
`

const PROMETHEUS_DELAY_METRICS_OUTPUT = `syslogng_output_event_delay_sample_seconds{transport="tcp",address="localhost:5555",driver="afsocket",id="#anon-destination0#0"} 5
syslogng_output_event_delay_sample_seconds{driver="http",url="http://localhost/asd",id="#anon-destination0#1",worker="0"} 2
syslogng_output_event_delay_sample_age_seconds{driver="http",url="http://localhost/asd",id="#anon-destination0#1",worker="0"} 1
syslogng_output_event_delay_sample_age_seconds{transport="tcp",address="localhost:5555",driver="afsocket",id="#anon-destination0#0"} 31
`

const PROMETHEUS_ESCAPE_METRICS_OUTPUT = `syslogng_classified_output_events_total{app="MSWinEventLog\\t1\\tSecurity\\t921448325\\tFri",source="s_critical_hosts_515"} 1
syslogng_classified_output_events_total{app="\a\t\n\"\xfa\\",source="s_unescaped_bug"} 1
`

func metricFamiliesToText(mfs []*io_prometheus_client.MetricFamily) string {
	var buf strings.Builder
	for _, mf := range mfs {
		_, _ = expfmt.MetricFamilyToText(&buf, mf)
	}
	return buf.String()
}

func sortMetricFamilies(mfs []*io_prometheus_client.MetricFamily) {
	slices.SortFunc(mfs, func(a, b *io_prometheus_client.MetricFamily) int {
		return strings.Compare(*a.Name, *b.Name)
	})
	for _, mf := range mfs {
		for _, m := range mf.Metric {
			slices.SortFunc(m.Label, func(a, b *io_prometheus_client.LabelPair) int {
				return strings.Compare(*a.Name, *b.Name)
			})
		}
		slices.SortFunc(mf.Metric, func(a, b *io_prometheus_client.Metric) int {
			return slices.CompareFunc(a.Label, b.Label, func(a, b *io_prometheus_client.LabelPair) int {
				if *a.Name != *b.Name {
					return strings.Compare(*a.Name, *b.Name)
				}
				return strings.Compare(*a.Value, *b.Value)
			})
		})
	}
}

func removeTimestamps(mfs []*io_prometheus_client.MetricFamily) {
	for _, mf := range mfs {
		for _, m := range mf.Metric {
			m.TimestampMs = nil
		}
	}
}

func amp[T any](v T) *T {
	return &v
}
