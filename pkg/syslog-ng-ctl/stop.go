// Copyright © 2025 Axoflow
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

import "context"

// Stop sends the stop command to the syslog-ng instance behind the control channel.
// With force set, axosyslog exits without waiting for its worker threads and loses the
// contents of its in-memory queues (syslog-ng-ctl stop --force, axoflow/axosyslog#1249).
// Older versions match control commands by prefix and run a plain STOP.
func Stop(ctx context.Context, cc ControlChannel, force bool) error {
	cmd := "STOP"
	if force {
		cmd = "STOP FORCE"
	}
	_, err := cc.SendCommand(ctx, cmd)
	return err
}
