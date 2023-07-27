// Copyright © 2023 Axoflow
// All rights reserved.

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	syslogngctl "github.com/axoflow/axoflow/go/pkg/syslog-ng-ctl"
	"github.com/prometheus/common/expfmt"
	"golang.org/x/exp/slices"
)

func main() {
	socketAddr := os.Getenv("CONTROL_SOCKET")
	if socketAddr == "" {
		_, _ = fmt.Fprintln(os.Stderr, "Control socket not specified. Set CONTROL_SOCKET environment variable.")
		os.Exit(1)
	}

	ctl := syslogngctl.Controller{
		ControlChannel: syslogngctl.NewUnixDomainSocketControlChannel(socketAddr),
	}

	cmds := []struct {
		Args []string
		Func func()
	}{
		{
			Args: []string{"ping"},
			Func: func() {
				err := ctl.Ping(context.Background())
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "An error occurred while pinging syslog-ng: %s\n", err.Error())
					os.Exit(2)
				}
			},
		},
		{
			Args: []string{"reload"},
			Func: func() {
				if err := ctl.Reload(context.Background()); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "An error occurred while reloading syslog-ng config: %s\n", err.Error())
					os.Exit(2)
				}
			},
		},
		{
			Args: []string{"show-license-info"},
			Func: func() {
				info, err := ctl.GetLicenseInfo(context.Background())
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "An error occurred while getting license info: %s\n", err.Error())
					os.Exit(2)
				}
				_, _ = fmt.Fprintln(os.Stdout, info)
			},
		},
		{
			Args: []string{"stats", "prometheus"},
			Func: func() {
				metrics, err := ctl.StatsPrometheus(context.Background())
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "An error occurred while querying prometheus stats: %s\n", err.Error())
					os.Exit(2)
				}
				for _, mf := range metrics {
					_, _ = expfmt.MetricFamilyToText(os.Stdout, mf)
				}
			},
		},
		{
			Args: []string{"stats"},
			Func: func() {
				stats, err := ctl.Stats(context.Background())
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "An error occurred while querying stats: %s\n", err.Error())
					os.Exit(2)
				}
				_, _ = fmt.Fprintf(os.Stdout, "%+v\n", stats)
			},
		},
	}

	for _, cmd := range cmds {
		if slices.Equal(os.Args[1:], cmd.Args) {
			cmd.Func()
			return
		}
	}
	_, _ = fmt.Fprintf(os.Stderr, "Unknown command %q\n", strings.Join(os.Args[1:], " "))
	_, _ = fmt.Fprintln(os.Stderr, "Supported commands:")
	for _, cmd := range cmds {
		_, _ = fmt.Fprintf(os.Stderr, "\t%s\n", strings.Join(cmd.Args, " "))
	}
	os.Exit(1)
}
