# Deploying an OpenTelemetry Collector

This demo OpenTelemetry collector was originally created for the [Export OpenTelemetry Data](https://docs.px.dev/tutorials/integrations/otel/) tutorial.

[OpenTelemetry Collectors](https://opentelemetry.io/docs/collector/) allow you to receive, process, and export telemetry data in a vendor-agnostic way. When configuring an OpenTelemetry Collector, you can specify how/where you should receive telemetry data, how that data should be processed, and how/where that data should be exported.

Our example OpenTelemetry Collector features a single Gateway cluster which receives data on the standard OpenTelemetry port (4317). This example collector outputs any data it receives to the logs.

## Instructions

1. Deploy the OpenTelemetry collector to your cluster using the `kubectl` command below. This will deploy the collector to the `default` namespace.

```
kubectl apply -f collector.yaml
```

2. [Configure the OpenTelemetry Pixie Plugin](https://docs.px.dev/tutorials/integrations/otel/#setup-the-plugin) to export data to `otel-collector.default.svc.cluster.local:4317`. Note that the OpenTelemetry collector must be deployed to the same cluster that Pixie is installed in.

3. To validate that data is being received by the OpenTelemetry collector, check the logs for the `otel-collector` pod. If the export was successful, you should see logs similar to:

```
2022-04-15T02:38:05.338Z	INFO	loggingexporter/logging_exporter.go:54	MetricsExporter	{"#metrics": 732}
```

## Have questions? Need help?

Please reach out on our Pixie Community [Slack](https://slackin.px.dev/) or file a GitHub issue.
