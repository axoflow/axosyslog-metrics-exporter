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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStopPicksTheControlCommand(t *testing.T) {
	for name, tc := range map[string]struct {
		force bool
		want  string
	}{
		"plain stop":  {want: "STOP"},
		"forced stop": {force: true, want: "STOP FORCE"},
	} {
		t.Run(name, func(t *testing.T) {
			var got string
			cc := ControlChannelFunc(func(_ context.Context, cmd string) (string, error) {
				got = cmd
				return "", nil
			})
			require.NoError(t, Stop(context.Background(), cc, tc.force))
			require.Equal(t, tc.want, got)
		})
	}
}
