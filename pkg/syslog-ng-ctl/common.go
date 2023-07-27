// Copyright © 2023 Axoflow
// All rights reserved.

package syslogngctl

import (
	"context"
	"fmt"
)

type ControlChannel interface {
	SendCommand(ctx context.Context, cmd string) (rsp string, err error)
}

type UnexpectedResponse string

func (err UnexpectedResponse) Error() string {
	return fmt.Sprintf("got unexpected response: %q", string(err))
}
