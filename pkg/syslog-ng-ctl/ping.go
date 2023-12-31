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
	"io"
)

// Ping checks whether syslog-ng is listening on the control channel
func Ping(ctx context.Context, cc ControlChannel) error {
	// send an inexpensive but valid command to prevent error logs
	rsp, err := License(ctx, cc)

	if err != nil && errors.Is(err, io.EOF) && rsp == "" {
		return nil // support very old syslog-ng versions that didn't support the LICENSE command
	}
	return err
}
