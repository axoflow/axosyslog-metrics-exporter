module github.com/axoflow/axosyslog-metrics-exporter

go 1.20

replace github.com/axoflow/axosyslog-metrics-exporter/pkg/syslog-ng-ctl => ./pkg/syslog-ng-ctl

require (
	github.com/axoflow/axosyslog-metrics-exporter/pkg/syslog-ng-ctl v0.0.0-00010101000000-000000000000
	github.com/prometheus/common v0.46.0
	golang.org/x/exp v0.0.0-20240119083558-1b970713d09a
)

require (
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/prometheus v0.50.1 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)
