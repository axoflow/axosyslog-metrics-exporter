module github.com/axoflow/axosyslog-metrics-exporter

go 1.24.11

require (
	github.com/axoflow/axosyslog-metrics-exporter/pkg/syslog-ng-ctl v0.0.0-20250721143838-ee0a5adf916c
	github.com/prometheus/common v0.67.2
	golang.org/x/exp v0.0.0-20250718183923-645b1fa84792
)

require (
	github.com/kr/text v0.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/prometheus v0.305.0 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace github.com/axoflow/axosyslog-metrics-exporter/pkg/syslog-ng-ctl => ./pkg/syslog-ng-ctl
