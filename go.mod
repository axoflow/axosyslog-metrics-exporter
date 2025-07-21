module github.com/axoflow/axosyslog-metrics-exporter

go 1.24.0

require (
	github.com/axoflow/axosyslog-metrics-exporter/pkg/syslog-ng-ctl v0.0.0-20250721143838-ee0a5adf916c
	github.com/prometheus/common v0.65.0
	golang.org/x/exp v0.0.0-20250718183923-645b1fa84792
)

require (
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/prometheus v0.305.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/axoflow/axosyslog-metrics-exporter/pkg/syslog-ng-ctl => ./pkg/syslog-ng-ctl
