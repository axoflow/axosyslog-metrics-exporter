// Copyright © 2023 Axoflow
// All rights reserved.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	syslogngctl "github.com/axoflow/axosyslog-metrics-exporter/pkg/syslog-ng-ctl"
	"github.com/prometheus/common/expfmt"
)

const (
	DEFAULT_TIMEOUT_SYSLOG time.Duration = time.Second * 5
	DEFAULT_SERVICE_PORT                 = "9577"
	DEFAULT_SOCKET_ADDR                  = "/var/run/syslog-ng/syslog-ng.tcl"
)

var (
	Version = "dev" // should be set build-time, see Makefile
)

type RunArgs struct {
	SocketAddr     string
	ServicePort    string
	RequestTimeout string
}

func envOrDef(envName string, def string) (res string) {
	res = os.Getenv(envName)
	if res == "" {
		res = def
	}
	return
}

func main() {
	runArgs := RunArgs{}

	fmt.Fprintf(os.Stdout, "%v version %q\n", filepath.Base(os.Args[0]), Version)

	flag.StringVar(&runArgs.SocketAddr, "socket.path", envOrDef("CONTROL_SOCKET", DEFAULT_SOCKET_ADDR), "syslog-ng control socket path")
	flag.StringVar(&runArgs.ServicePort, "service.port", envOrDef("SERVICE_PORT", DEFAULT_SERVICE_PORT), "service port")
	flag.StringVar(&runArgs.RequestTimeout, "service.timeout", envOrDef("SERVICE_TIMEOUT", DEFAULT_TIMEOUT_SYSLOG.String()), "request timeout")

	flag.Parse()

	requestTimeout, err := time.ParseDuration(runArgs.RequestTimeout)
	if err != nil {
		requestTimeout = DEFAULT_TIMEOUT_SYSLOG
	}

	ctl := syslogngctl.Controller{
		ControlChannel: syslogngctl.NewUnixDomainSocketControlChannel(runArgs.SocketAddr),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {

		subCtx, cancel := context.WithTimeout(r.Context(), requestTimeout)
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
		subCtx, cancel := context.WithTimeout(r.Context(), requestTimeout)
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

	fmt.Fprintf(os.Stdout, "listening on port: %v\n", runArgs.ServicePort)
	fmt.Fprintf(os.Stdout, "syslog-ng control socket path: %v\n", runArgs.SocketAddr)
	fmt.Fprintf(os.Stdout, "service timeout: %v\n", requestTimeout.String())

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%v", runArgs.ServicePort), mux))
}
