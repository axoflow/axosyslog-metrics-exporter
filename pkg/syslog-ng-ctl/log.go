// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import "strings"

func getLog(cc ControlChannel, mode string) (on bool, err error) {
	rsp, err := cc.SendCommand("LOG " + mode)
	if err != nil {
		return
	}
	for _, line := range strings.Fields(rsp) {
		switch line {
		case "OK": // 3.35.1: OK verbose=0
		case "syslog-ng", "log", "level", "set", "to": //  4.1.1: OK syslog-ng log level set to 0
		case mode + "=0", "0":
			on = false
			return
		case mode + "=1", "1":
			on = true
			return
		default:
			err = UnexpectedResponse(rsp)
			return
		}
	}

	err = UnexpectedResponse(rsp)
	return
}

func setLog(cc ControlChannel, mode string, on bool) error {
	cmd := "LOG " + mode
	if on {
		cmd += " ON"
	} else {
		cmd += " OFF"
	}
	_, err := cc.SendCommand(cmd)
	return err
}
