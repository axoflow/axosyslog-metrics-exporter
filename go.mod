module github.com/axoflow/axosyslog-metrics-exporter

go 1.25.0

require (
	github.com/axoflow/axosyslog-metrics-exporter/pkg/syslog-ng-ctl v0.0.0-20250721143838-ee0a5adf916c
	github.com/prometheus/common v0.69.0
	golang.org/x/exp v0.0.0-20260218203240-3dfff04db8fa
)

require (
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/prometheus v0.313.1 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/axoflow/axosyslog-metrics-exporter/pkg/syslog-ng-ctl => ./pkg/syslog-ng-ctl
