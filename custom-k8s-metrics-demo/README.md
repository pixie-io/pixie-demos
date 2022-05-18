# Pixie Custom Kubernetes Metrics Demo

Autoscale the number of pods in your Kubernetes deployment based on request throughput, without any code changes. [Horizontal Pod Autoscaling with Custom Metrics in Kubernetes](https://blog.px.dev/autoscaling-custom-k8s-metric) is the accompanying blog post for this demo.

## What is this demo?

This demo provides an example implementation of a custom Kubernetes metric from [Pixie](https://github.com/pixie-io/pixie) data. Specifically, it provides a metric for the number of HTTP requests per second by Kubernetes pod. This custom metric can then be used as an input to a [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) in order to to scale up/down the number of pods.

This demo was based off of the example in [kubernetes-sigs/custom-metrics-apiserver](https://github.com/kubernetes-sigs/custom-metrics-apiserver).

You can also view a live version of this demo at this [talk](https://www.youtube.com/watch?v=EG4isSqD3IE) (autoscaling content at about 22:43 in).

## Prerequisites

* If you don't already have one, set up a [Kubernetes cluster](https://docs.px.dev/installing-pixie/setting-up-k8s/)
* [Install Pixie](https://docs.px.dev/installing-pixie/install-guides/) to your Kubernetes cluster
* [Optional] Install [hey](https://github.com/rakyll/hey) for testing the demo application 

## Usage

1. Create a secret containing the Pixie API credentials for your Kubernetes cluster:

```
# Create `px-custom-metrics` namespace
kubectl create namespace px-custom-metrics

# Get your current cluster name from your Kubernetes context
kubectl config current-context
# Get the Pixie Cluster ID for the above cluster name.
# Record the value of the `ID` column for this cluster
px get viziers

# Create an API key
# Record the value of the `Key` parameter
px api-key create

kubectl -n px-custom-metrics create secret generic px-credentials --from-literal=px-api-key=<YOUR API KEY VALUE HERE> --from-literal=px-cluster-id=<YOUR CLUSTER ID VALUE HERE>
```

2. [Optional] If using self-hosted Pixie Cloud, update PX_CLOUD_ADDR in `px-custom-metrics.yaml`.

3. Create the Pixie metrics provider in your Kubernetes cluster in the `px-custom-metrics` namespace:

```
kubectl apply -f px-custom-metrics.yaml
```
4. Wait until the pods in the `px-custom-metrics` namespace are up and healthy.

5. Check to make sure that the metric server returns metrics as expected:

```
kubectl -n px-custom-metrics get --raw "/apis/custom.metrics.k8s.io/v1beta2/namespaces/default/pods/*/px-http-requests-per-second"
```

## Test application

1. Deploy a test application to autoscale based on the metrics you just created.

```
# Will be created in the `default` namespace
kubectl apply -f demo-app.yaml
```

2. Check the number of pods that are currently running for `default/echo-service` (should be 1).

```
kubectl get pods --selector=name=echo-service
```

3. Get the external IP for `echo-service`:

```
kubectl get services
```

4. Increase the load on `echo-service` with the external IP from the previous step:

```
hey -n 10000  http://<EXTERNAL IP>/ping
```

5. Watch the number of pods for `echo-service`. Within a few minutes, it should go up, and then go back down to 1 pod after a while.

```
kubectl get pods --selector=name=echo-service --watch
```

6. Test out different endpoints in the echo-service.

* `/ping` will echo back the body of the request
* `/expensive` will do the same as `/ping`, but add an expensive-ish calculation to utilize CPU
* `/slow-contention` will add artificial delay in responses for queued requests in order to simulate a non-CPU bottleneck such as another service
* `/expensive-limit-concurrent` will make an expensive computation but return errors once a certain number of concurrent requests occur, in order to simulate a request queue with a limit.

These requests can be tested in conjunction with the different metrics specified in `pixie-http-metric-provider.go`. Just edit the metric name in the HorizontalPodAutoscaler from the default of `px-http-requests-per-second`.


## Development

You can use the pre-built images in the directions described above if you don't need to make any changes to the metrics provider code. However, if you want to extend or develop upon this example, here are instructions for building and deploying a local version of the metrics server.

1. Make your changes and build a new version of the server image:

```
docker build . -t <YOUR DOCKER IMAGE PATH HERE> 
```

2. Push your version of the image (not necessary if your cluster has access to your local Docker images):

```
docker push <YOUR DOCKER IMAGE PATH HERE>:latest
```

3. Depending on your ImagePullPolicy, delete and recreate the `px-custom-metrics` deployment:

```
kubectl -n px-custom-metrics delete deployment px-custom-metrics-apiserver
kubectl apply -f px-custom-metrics.yaml
```

## Extensions
This is an example implementation of a Pixie custom metrics server. Pixie can be used to generate many different types of metrics, not just HTTP request throughput by pod.

Pixie can generate metrics by pod, service, node, or container.

Other example metrics Pixie can generate:
* Latency, error rate, and throughput for our [supported protocols](https://docs.px.dev/about-pixie/data-sources/#supported-protocols).
* Latency, error rate, and throughput by request path (including wildcards, such as `/orders/*/item/*`)
* System metrics such as CPU, network utilization, memory utilization
* Application CPU profiles
* See our example [PxL scripts](https://github.com/pixie-io/pixie/tree/main/src/pxl_scripts) for additional examples

## Bugs & Features

Feel free to file a bug or an issue for a feature request. You can also join our [Slack](https://slackin.px.dev/) community.
