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
	DEFAULT_SOCKET_ADDR                  = "/var/run/syslog-ng/syslog-ng.ctl"
)

var (
	Version = "dev" // should be set build-time, see Makefile
)

type RunArgs struct {
	SocketAddr     string
	ServicePort    string
	ServiceAddress string
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
	flag.StringVar(&runArgs.ServicePort, "service.port", envOrDef("SERVICE_PORT", DEFAULT_SERVICE_PORT), "service bind port")
	flag.StringVar(&runArgs.ServiceAddress, "service.address", envOrDef("SERVICE_ADDRESS", ""), "service bind address in [host]:port format (overwrites service.port)")
	flag.StringVar(&runArgs.RequestTimeout, "service.timeout", envOrDef("SERVICE_TIMEOUT", DEFAULT_TIMEOUT_SYSLOG.String()), "request timeout")

	flag.Parse()
	if runArgs.ServiceAddress == "" {
		runArgs.ServiceAddress = fmt.Sprintf(":%v", runArgs.ServicePort)
	}

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

	err = http.ListenAndServe(runArgs.ServiceAddress, mux)
	fmt.Println(err)
}
