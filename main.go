// Copyright © 2023 Axoflow
// All rights reserved.

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	syslogngctl "github.com/axoflow/metrics-exporter/pkg/syslog-ng-ctl"
	"github.com/prometheus/common/expfmt"
)

const (
	TIMEOUT_SYSLOG time.Duration = time.Second * 3
	HTTP_PORT                    = 9999
)

var (
	Version = "dev" // should be set build-time, see Makefile
)

func main() {
	fmt.Fprintf(os.Stdout, "metrics exporter version %q\n", Version)

	socketAddr := os.Getenv("CONTROL_SOCKET")
	if socketAddr == "" {
		socketAddr = "/var/run/syslog-ng/syslog-ng.tcl"
		_, _ = fmt.Fprintf(os.Stdout, "CONTROL_SOCKET environment variable not set, defaulting to %q", socketAddr)
	}

	ctl := syslogngctl.Controller{
		ControlChannel: syslogngctl.NewUnixDomainSocketControlChannel(socketAddr),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		subCtx, cancel := context.WithTimeout(r.Context(), TIMEOUT_SYSLOG)
		defer cancel()
		mfs, err := ctl.StatsPrometheus(subCtx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var resp bytes.Buffer

		for _, mf := range mfs {
			_, err := expfmt.MetricFamilyToText(&resp, mf)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		_, err = io.Copy(w, &resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		subCtx, cancel := context.WithTimeout(r.Context(), TIMEOUT_SYSLOG)
		defer cancel()
		err := ctl.Ping(subCtx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write([]byte(`PONG`))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Fprintf(os.Stdout, "listening on port %v, \n", HTTP_PORT)
	fmt.Fprintf(os.Stdout, "syslog-ng control socket path %v, \n", socketAddr)

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%v", HTTP_PORT), mux))
}
