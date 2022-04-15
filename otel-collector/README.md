# Deploying an OpenTelemetry Collector

[OpenTelemetry Collectors](https://opentelemetry.io/docs/collector/) allow you to receive, process, and export telemetry data in a vendor-agnostic way. When configuring an OpenTelemetry Collector, you can specify how/where you should receive telemetry data, how that data should be processed, and how/where that data should be exported.

Our example OpenTelemetry Collector features a single Gateway cluster which receives data on the standard OpenTelemetry port (4317). This example collector outputs any data it receives to the logs.

To deploy:

```
kubectl apply -f collector.yaml
```

This will deploy the collector to the `default` namespace. Send data to `otel-collector.default.svc.cluster.local:4317` from within the cluster to export data to the collector.
To validate that the data has been received, simply view the `otel-collector` pod's logs. If the export was successful, you should see logs similar to:

```
2022-04-15T02:38:05.338Z	INFO	loggingexporter/logging_exporter.go:54	MetricsExporter	{"#metrics": 732}
```
