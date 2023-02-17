package syslogngctl

import (
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/multierr"

	"github.com/axoflow/axo-edge/internal/syslog-ng/stats"
)

// Controller implements syslog-ng-ctl's functionality.
//
// Reference for available commands in syslog-ng-ctl's source code: https://github.com/syslog-ng/syslog-ng/blob/0e7c762c704efbda0ae10b61c35700ef0bdbb9c1/syslog-ng-ctl/syslog-ng-ctl.c#L111
type Controller struct {
	Runner CommandRunner
}

func (c Controller) GetDebug() (on bool, err error) {
	return c.getLog("DEBUG")
}

func (c Controller) GetTrace() (on bool, err error) {
	return c.getLog("TRACE")
}

func (c Controller) GetVerbose() (on bool, err error) {
	return c.getLog("VERBOSE")
}

func (c Controller) Reload() error {
	_, err := c.runCmd("RELOAD")
	return err
}

func (c Controller) SetDebug(on bool) error {
	return c.setLog("DEBUG", on)
}

func (c Controller) SetTrace(on bool) error {
	return c.setLog("TRACE", on)
}

func (c Controller) SetVerbose(on bool) error {
	return c.setLog("VERBOSE", on)
}

func (c Controller) Stats() (res []stats.Stat, errs error) {
	rsp, err := c.runCmd("STATS")
	if err != nil {
		return res, err
	}
	rsp = strings.TrimRight(rsp, "\n") // remove trailing new line
	lines := strings.Split(rsp, "\n")
	// TODO: sanity check: match header
	lines = lines[1:] // drop header line: SourceName;SourceId;SourceInstance;State;Type;Number
	for _, line := range lines {
		fields := strings.Split(line, ";")
		if len(fields) != 6 {
			errs = multierr.Append(errs, InvalidStatLine(line))
			continue
		}
		if len(fields[3]) != 1 {
			errs = multierr.Append(errs, InvalidStatLine(line))
			continue
		}
		num, err := strconv.ParseUint(fields[5], 10, 64)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
		res = append(res, stats.Stat{
			SourceName:     fields[0],
			SourceID:       fields[1],
			SourceInstance: fields[2],
			SourceState:    stats.SourceState(fields[3][0]),
			Type:           fields[4],
			Number:         num,
		})
	}
	return
}

func (c Controller) getLog(mode string) (on bool, err error) {
	rsp, err := c.runCmd("LOG " + mode)
	if err != nil {
		return
	}
	switch strings.TrimRight(rsp, "\n") {
	case mode + "=0":
		on = false
	case mode + "=1":
		on = true
	default:
		err = UnexpectedResponse(rsp)
	}
	return
}

func (c Controller) runCmd(cmd string) (string, error) {
	return c.Runner.RunCommand(cmd)
}

func (c Controller) setLog(mode string, on bool) error {
	cmd := "LOG " + mode
	if on {
		cmd += " ON"
	} else {
		cmd += " OFF"
	}
	_, err := c.runCmd(cmd)
	return err
}

type CommandRunner interface {
	RunCommand(cmd string) (rsp string, err error)
}

type UnexpectedResponse string

func (err UnexpectedResponse) Error() string {
	return fmt.Sprintf("got unexpected response: %q", string(err))
}

type InvalidStatLine string

func (err InvalidStatLine) Error() string {
	return fmt.Sprintf("invalid stat line: %q", string(err))
}
