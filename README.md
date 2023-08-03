<p align="center">
  <picture>
    <source media="(prefers-color-scheme: light)" srcset="https://github.com/axoflow/axosyslog-docker/raw/main/docs/axoflow-logo-color.svg">
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/axoflow/axosyslog-docker/raw/main/docs/axoflow-logo-white.svg">
    <img alt="Axoflow" src="https://github.com/axoflow/axosyslog-docker/raw/main/docs/axoflow-logo-color.svg" width="550">
  </picture>
</p>

# axosyslog-metrics-exporter

Export [prometheus stats](https://axoflow.com/docs/axosyslog/docs/parsers/metrics-probe/) of Axosyslog over HTTP.

## About

Axosyslog-metrics-exporter serves Prometheus metrics over a HTTP interface (`http://0.0.0.0:9577/metrics` by default).
It needs UNIX file-level access to syslog-ng's control socket, which is usually at
`/var/lib/syslog-ng/syslog-ng.ctl` or `/var/run/syslog-ng/syslog-ng.ctl`).
In container environments you need to provide access to that UNIX domain socket via shared volumes or other means.

The HTTP and command line interface is compatible with [syslog_ng_exporter](https://github.com/kube-logging/syslog_ng_exporter),
but we use the new native prometheus stats available since versions 4.1.
We keep translating from the legacy `stats` interface in case of older syslog-ng versions.

## Usage

### Command line

```
axosyslog-metrics-exporter [options]

Options:
  -service.port string
    	service bind port (default "9577" or $SERVICE_PORT)
  -service.timeout string
    	request timeout (default "5s" or $SERVICE_TIMEOUT)
  -socket.path string
    	syslog-ng control socket path (default "/var/run/syslog-ng/syslog-ng.ctl" or $CONTROL_SOCKET)
```

### Docker

```
docker run -d -p 9577:9577 -v $(echo /var/*/syslog-ng/syslog-ng.ctl):/syslog-ng.ctl \
  ghcr.io/axoflow/axosyslog-metrics-exporter:latest --socket.path=/syslog-ng.ctl
```

### Logging-operator

You can replace the `exporter` sidecar's image in syslog-ng based logging-operator setups, by extending the Logging resource:

```
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: my-logging
spec:
  controlNamespace: logging-operator
  loggingRef: my-logging
  syslogNG:
    globalOptions:
      stats:
        freq: 0
        level: 2
    statefulSet:
      spec:
        template:
          spec:
            containers:
            - image: ghcr.io/axoflow/axosyslog-metrics-exporter:latest
              name: exporter
```

## Contact and support

In case you need help or want to contact us, open a [GitHub issue](https://github.com/axoflow/axosyslog-metrics-exporter/issues), or come chat with us in the [syslog-ng channel of the Axoflow Discord server](https://discord.gg/4Fzy7D66Qq).

## Contribution

If you have fixed a bug or would like to contribute your improvements to these images, [open a pull request](https://github.com/axoflow/axosyslog-metrics-exporter/pulls). We truly appreciate your help.

## About Axoflow

The [Axoflow](https://axoflow.com) founder team consists of successful entrepreneurs with a vast knowledge and hands-on experience about observability, log management, and how to apply these technologies in the enterprise security context. We also happen to be the creators of wide-spread open source technologies in this area, like syslog-ng and the [Logging operator for Kubernetes](https://github.com/kube-logging/logging-operator).

To learn more about our products and our open-source projects, visit the [Axoflow blog](https://axoflow.com/blog/), or [subscribe to the Axoflow newsletter](https://axoflow.com/#newsletter-subscription).
