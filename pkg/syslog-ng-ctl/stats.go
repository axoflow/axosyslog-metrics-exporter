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
	"fmt"
	"strconv"
	"strings"
)

func Stats(ctx context.Context, cc ControlChannel) ([]Stat, error) {
	rsp, err := cc.SendCommand(ctx, "STATS")
	if err != nil {
		return nil, err
	}

	return parseStats(rsp)
}

type Stat struct {
	SourceName     string
	SourceID       string
	SourceInstance string
	SourceState    SourceState
	Type           string
	Number         uint64
}

type SourceState byte

const (
	SourceStateActive   SourceState = 'a'
	SourceStateDynamic  SourceState = 'd'
	SourceStateOrphaned SourceState = 'o'
)

type InvalidStatLine string

func (err InvalidStatLine) Error() string {
	return fmt.Sprintf("invalid stat line: %q", string(err))
}

const StatsHeader = "SourceName;SourceId;SourceInstance;State;Type;Number"

func parseStats(rsp string) (stats []Stat, errs error) {
	rsp = strings.TrimRight(rsp, "\n") // remove trailing new line
	lines := strings.Split(rsp, "\n")
	// TODO: sanity check: match header line
	lines = lines[1:] // drop header line: SourceName;SourceId;SourceInstance;State;Type;Number
	for _, line := range lines {
		fields := strings.Split(line, ";")
		if len(fields) != 6 {
			errs = errors.Join(errs, InvalidStatLine(line))
			continue
		}
		if len(fields[3]) != 1 {
			errs = errors.Join(errs, InvalidStatLine(line))
			continue
		}
		num, err := strconv.ParseUint(fields[5], 10, 64)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		if fields[3] != string(SourceStateOrphaned) {
			stats = append(stats, Stat{
				SourceName:     fields[0],
				SourceID:       fields[1],
				SourceInstance: fields[2],
				SourceState:    SourceState(fields[3][0]),
				Type:           fields[4],
				Number:         num,
			})
		}
	}
	return
}
