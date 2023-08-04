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
)

// GetLicenseInfo sends the LICENSE command to syslog-ng
func GetLicenseInfo(ctx context.Context, cc ControlChannel) (string, error) {
	info, err := cc.SendCommand(ctx, "LICENSE")
	info = strings.TrimSpace(info)
	return info, err
}
