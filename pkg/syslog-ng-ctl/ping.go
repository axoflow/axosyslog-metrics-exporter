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
	"strings"
)

// Ping checks whether syslog-ng is listening on the control channel
func Ping(ctx context.Context, cc ControlChannel) error {
	rsp, err := cc.SendCommand(ctx, "invalid") // we just send a command to check if syslog-ng is listening on the other end
	if strings.TrimSpace(rsp) == "" && errors.Is(err, io.EOF) {
		return nil
	}
	return err
}
