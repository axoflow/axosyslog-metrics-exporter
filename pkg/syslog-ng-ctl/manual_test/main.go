package main

import (
	"fmt"
	"io"
	"net"
	"os"

	syslogngctl "github.com/axoflow/axo-edge/pkg/syslog-ng-ctl"
)

func main() {
	socketAddr := os.Getenv("CONTROL_SOCKET")
	if socketAddr == "" {
		_, _ = fmt.Fprintln(os.Stderr, "Control socket not specified. Set CONTROL_SOCKET environment variable.")
		os.Exit(1)
	}
	ctl := syslogngctl.Controller{
		ControlChannel: syslogngctl.NewReadWriterControlChannel(func() (io.ReadWriter, error) {
			conn, err := net.Dial("unix", socketAddr)
			return conn, err
		}),
	}
	stats, err := ctl.Stats()
	_, _ = fmt.Fprintf(os.Stdout, "%+v\n", stats)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "An error occurred while querying stats: %s\n", err.Error())
		os.Exit(1)
	}
}
