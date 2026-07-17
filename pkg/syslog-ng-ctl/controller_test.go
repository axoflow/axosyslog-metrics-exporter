// Copyright © 2026 Axoflow
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
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestControllerStatsPrometheusConcurrent(t *testing.T) {
	ctl := NewController(ControlChannelFunc(func(_ context.Context, cmd string) (string, error) {
		require.Equal(t, "STATS PROMETHEUS", cmd)
		return "syslogng_output_event_delay_sample_seconds{id=\"d\"} 1\n" +
			"syslogng_output_event_delay_sample_age_seconds{id=\"d\"} 0\n", nil
	}))

	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := ctl.StatsPrometheus(context.Background())
			require.NoError(t, err)
		}()
	}
	wg.Wait()
}
