module github.com/axoflow/axosyslog-metrics-exporter

go 1.23.5

replace github.com/axoflow/axosyslog-metrics-exporter/pkg/syslog-ng-ctl => ./pkg/syslog-ng-ctl

require (
	github.com/axoflow/axosyslog-metrics-exporter/pkg/syslog-ng-ctl v0.0.0-20241004091155-72f25f49e310
	github.com/prometheus/common v0.62.0
	golang.org/x/exp v0.0.0-20250128182459-e0ece0dbea4c
)

require (
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/prometheus v0.301.0 // indirect
	google.golang.org/protobuf v1.36.4 // indirect
)
