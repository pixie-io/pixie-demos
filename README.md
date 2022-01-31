# pixie-demos

## [custom-k8s-metrics-demo](https://github.com/pixie-io/pixie-demos/tree/main/custom-k8s-metrics-demo)

This demo project shows how to use Pixie to autoscale the number of pods in your Kubernetes deployment based on request throughput, without any code changes. [Horizontal Pod Autoscaling with Custom Metrics in Kubernetes](https://blog.px.dev/autoscaling-custom-k8s-metric) is the accompanying blog post.

## [ebpf-profiler](https://github.com/pixie-io/pixie-demos/tree/main/ebpf-profiler)

Demo project to accompany the [Building a Continuous Profiler Part 2: A Simple eBPF-Based Profiler](https://blog.px.dev/cpu-profiling-2/) blog post. This CPU performance profiler project shows how to get sample stack traces for performance profiling, using eBPF.

## [eks-workshop](https://github.com/pixie-io/pixie-demos/tree/main/eks-workshop)

Resources for the [Monitoring with Pixie](https://www.eksworkshop.com/intermediate/241_pixie/) AWS EKS Workshop.

## [endpoint-deprecation](https://github.com/pixie-io/pixie-demos/tree/main/endpoint-deprecation)

Want to deprecate an API? Use [Pixie](https://github.com/pixie-io/pixie) to quickly determine:

- Is this endpoint used?
- Who is using this endpoint?

[Can I deprecate this endpoint?](https://blog.px.dev/endpoint-deprecation) is the accompanying blog post for this demo.

## [http2-tracing](https://github.com/pixie-io/pixie-demos/tree/main/http2-tracing)

Demo project to accompany the [Observing HTTP/2 Traffic is Hard, but eBPF Can Help](https://blog.px.dev/ebpf-http2-tracing/) blog post. This is a basic example of how to trace HTTP/2 messages using eBPF uprobes.

## [k8s-cost-estimation](https://github.com/pixie-io/pixie-demos/tree/main/k8s-cost-estimation)

Use Pixie to estimate the cost of hosting your Kubernetes cluster.

## [openssl-tracer](https://github.com/pixie-io/pixie-demos/tree/main/openssl-tracer)

Demo project to accompany the [Debugging with eBPF Part 3: Tracing SSL/TLS connections](https://blog.px.dev/ebpf-openssl-tracing/) blog post. This is a basic example of how to trace the OpenSSL library using eBPF. This tracer uses BCC to deploy the eBPF probes.

## [react-table](https://github.com/pixie-io/pixie-demos/tree/main/react-table)

Demo project to accompany the [Tables are Hard, Part 2: Building a Simple Data Table in React](https://blog.px.dev/tables-are-hard-2) blog post. Interactive demo: [github.io](https://pixie-io.github.io/pixie-demos/react-table).

## [simple-gotracing](https://github.com/pixie-io/pixie-demos/tree/main/simple-gotracing)

Demo project to accompany the [Dynamic Logging in Go](https://docs.pixielabs.ai/tutorials/custom-data/dynamic-go-logging/) tutorial.

## [slack-alert-app](https://github.com/pixie-io/pixie-demos/tree/main/slack-alert-app)

Demo project to accompany the [Slack Alerts using the Pixie API](https://docs.pixielabs.ai/tutorials/integrations/slackbot-alert/) tutorial. This demo project creates a Slackbot that reports the number of HTTP errors per service in your cluster.

## [sql-injection-demo](https://github.com/pixie-io/pixie-demos/tree/main/sql-injection-demo)

Demo project to accompany the [Detect SQL injections with Pixie](https://blog.px.dev/sql-injection/) blog post. This demo shows how to use Pixie to detect SQL injections on a Kubernetes application.

# Have questions? Need help?

Please reach out on our Pixie Community [Slack](https://slackin.px.dev/) or file a GitHub issue.
