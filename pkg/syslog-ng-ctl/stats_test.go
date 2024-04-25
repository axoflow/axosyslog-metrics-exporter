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
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStats(t *testing.T) {
	response := `SourceName;SourceId;SourceInstance;State;Type;Number
global;payload_reallocs;;a;processed;6
destination;d_cron;;a;processed;0
global;internal_queue_length;;a;processed;0
destination;d_mail;;a;processed;0
destination;d_console_all;;a;processed;0
destination;d_xconsole;;a;processed;0
destination;d_kern;;a;processed;0
source;s_src;;a;processed;65
global;sdata_updates;;a;processed;0
destination;d_console;;a;processed;0
destination;d_error;;a;processed;0
destination;d_daemon;;a;processed;0
center;;queued;a;processed;124
destination;d_debug;;a;processed;0
destination;d_uucp;;a;processed;0
destination;d_messages;;a;processed;59
destination;d_newscrit;;a;processed;0
global;msg_clones;;a;processed;6
global;scratch_buffers_count;;a;queued;8589934592
destination;d_user;;a;processed;0
destination;d_syslog;;a;processed;59
center;;received;a;processed;65
destination;d_newsnotice;;a;processed;0
destination;d_auth;;a;processed;6
destination;d_lpr;;a;processed;0
src.internal;s_src#1;;a;processed;59
src.internal;s_src#1;;a;stamp;1673105444
destination;d_newserr;;a;processed;0
global;scratch_buffers_bytes;;a;queued;0
dst.network;#anon-destination0#0;tcp,localhost:1234;o;eps_last_1h;0
dst.network;#anon-destination0#0;tcp,localhost:1234;o;eps_last_24h;0
dst.network;#anon-destination0#0;tcp,localhost:1234;o;dropped;0
dst.network;#anon-destination0#0;tcp,localhost:1234;o;processed;0
dst.network;#anon-destination0#0;tcp,localhost:1234;o;queued;0
dst.network;#anon-destination0#0;tcp,localhost:1234;o;memory_usage;0
dst.network;#anon-destination0#0;tcp,localhost:1234;o;written;0
dst.network;#anon-destination0#0;tcp,localhost:1234;o;truncated_bytes;0
dst.network;#anon-destination0#0;tcp,localhost:1234;o;eps_since_start;0
dst.network;#anon-destination0#0;tcp,localhost:1234;o;msg_size_max;0
dst.network;#anon-destination0#0;tcp,localhost:1234;o;truncated_count;0
dst.network;#anon-destination0#0;tcp,localhost:1234;o;msg_size_avg;0
.
`
	expected := []Stat{
		{SourceName: "global", SourceID: "payload_reallocs", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 6},
		{SourceName: "destination", SourceID: "d_cron", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "global", SourceID: "internal_queue_length", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "destination", SourceID: "d_mail", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "destination", SourceID: "d_console_all", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "destination", SourceID: "d_xconsole", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "destination", SourceID: "d_kern", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "source", SourceID: "s_src", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 65},
		{SourceName: "global", SourceID: "sdata_updates", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "destination", SourceID: "d_console", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "destination", SourceID: "d_error", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "destination", SourceID: "d_daemon", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "center", SourceID: "", SourceInstance: "queued", SourceState: SourceStateActive, Type: "processed", Number: 124},
		{SourceName: "destination", SourceID: "d_debug", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "destination", SourceID: "d_uucp", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "destination", SourceID: "d_messages", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 59},
		{SourceName: "destination", SourceID: "d_newscrit", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "global", SourceID: "msg_clones", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 6},
		{SourceName: "global", SourceID: "scratch_buffers_count", SourceInstance: "", SourceState: SourceStateActive, Type: "queued", Number: 8589934592},
		{SourceName: "destination", SourceID: "d_user", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "destination", SourceID: "d_syslog", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 59},
		{SourceName: "center", SourceID: "", SourceInstance: "received", SourceState: SourceStateActive, Type: "processed", Number: 65},
		{SourceName: "destination", SourceID: "d_newsnotice", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "destination", SourceID: "d_auth", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 6},
		{SourceName: "destination", SourceID: "d_lpr", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "src.internal", SourceID: "s_src#1", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 59},
		{SourceName: "src.internal", SourceID: "s_src#1", SourceInstance: "", SourceState: SourceStateActive, Type: "stamp", Number: 1673105444},
		{SourceName: "destination", SourceID: "d_newserr", SourceInstance: "", SourceState: SourceStateActive, Type: "processed", Number: 0},
		{SourceName: "global", SourceID: "scratch_buffers_bytes", SourceInstance: "", SourceState: SourceStateActive, Type: "queued", Number: 0},
	}
	request := bytes.Buffer{}
	cc := NewReadWriterControlChannel(func(context.Context) (io.ReadWriter, error) {
		return struct {
			io.Reader
			io.Writer
		}{
			Reader: strings.NewReader(response),
			Writer: &request,
		}, nil
	})
	res, err := Stats(context.Background(), cc)
	require.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.Equal(t, "STATS\n", request.String())
}
