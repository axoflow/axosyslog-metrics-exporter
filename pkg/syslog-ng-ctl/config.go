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

import "context"

// OriginalConfig sends the CONFIG GET ORIGINAL command to syslog-ng
func OriginalConfig(ctx context.Context, cc ControlChannel) (string, error) {
	return cc.SendCommand(ctx, "CONFIG GET ORIGINAL")
}

// PreprocessedConfig sends the CONFIG GET PREPROCESSED command to syslog-ng
func PreprocessedConfig(ctx context.Context, cc ControlChannel) (string, error) {
	return cc.SendCommand(ctx, "CONFIG GET PREPROCESSED")
}

// PreprocessedConfig sends the CONFIG ID command to syslog-ng
func ConfigID(ctx context.Context, cc ControlChannel) (string, error) {
	return cc.SendCommand(ctx, "CONFIG ID")
}
