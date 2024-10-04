module github.com/axoflow/axosyslog-metrics-exporter

go 1.22.0

toolchain go1.23.1

replace github.com/axoflow/axosyslog-metrics-exporter/pkg/syslog-ng-ctl => ./pkg/syslog-ng-ctl

require (
	github.com/axoflow/axosyslog-metrics-exporter/pkg/syslog-ng-ctl v0.0.0-20240731091211-4160b5cc192f
	github.com/prometheus/common v0.60.0
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0
)

require (
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/prometheus v0.54.1 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
