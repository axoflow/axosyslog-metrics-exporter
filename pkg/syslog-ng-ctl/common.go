// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import "fmt"

type ControlChannel interface {
	SendCommand(cmd string) (rsp string, err error)
}

type UnexpectedResponse string

func (err UnexpectedResponse) Error() string {
	return fmt.Sprintf("got unexpected response: %q", string(err))
}
