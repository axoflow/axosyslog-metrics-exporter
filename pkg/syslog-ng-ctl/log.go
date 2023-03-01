// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import "strings"

func getLog(cc ControlChannel, mode string) (on bool, err error) {
	rsp, err := cc.SendCommand("LOG " + mode)
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
