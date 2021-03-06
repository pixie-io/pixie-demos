# Observability for Feature Deprecation

Want to deprecate an API? Use [Pixie](https://github.com/pixie-io/pixie) to quickly determine:

- Is this endpoint used?
- Who is using this endpoint?

[Can I deprecate this endpoint?]( https://blog.px.dev/endpoint-deprecation) is the accompanying blog post for this demo.

## Prerequisites

- If you don't already have one, set up a [Kubernetes cluster](https://docs.px.dev/installing-pixie/setting-up-k8s/).
- [Install Pixie](https://docs.px.dev/installing-pixie/install-guides/) to your Kubernetes cluster.
- Install the [Pixie CLI](https://docs.px.dev/installing-pixie/install-schemes/cli/#1.-install-the-pixie-cli) if you didn't install it in order to deploy Pixie in the step above.
- git clone this repository and `cd` into the `endpoint-deprecation` folder.

## Test Application

1. Install [hey](https://github.com/rakyll/hey). Homebrew users can use:

```
brew install hey
```

2. Deploy an echo server to receive arbitrary client requests:

```
# echo-service will be created in the `default` namespace.
kubectl apply -f demo-app.yaml
```

3. Get the external IP for `echo-service` and save it in an environment variable:

```
kubectl get services
export ECHO_SERVICE_IP=<EXTERNAL IP>
```

4. Run the test load. `-H` is used to pass a custom HTTP header.

```
for i in {1..15}; do
hey -H "Referer:https://example.com/" -H "API-KEY:abcdef12345" -n 550 "http://${ECHO_SERVICE_IP}/v1/catalog/"
hey -H "Referer:https://px.dev/" -H "API-KEY:lkjlsdfsdfs" -n 50 "http://${ECHO_SERVICE_IP}/v1/catalog/$(uuidgen)/details"
hey -H "Referer:https://example.com/" -H "API-KEY:abcdef12345" -n 50 "http://${ECHO_SERVICE_IP}/v1/catalog/$(uuidgen)/details"
hey -H "Referer:https://example.com/" -H "API-KEY:sdfsdfsdfsd" -n 50 "http://${ECHO_SERVICE_IP}/v1/catalog/$(uuidgen)/details"
hey -H "Referer:https://google.com/" -H "API-KEY:sdfsdfsdfsd" -n 50 "http://${ECHO_SERVICE_IP}/v2/catalog/$(uuidgen)"
done
```

## Service Traffic Clustered by Logical Endpoint

From the top-level `endpoint-deprecation` folder, run:

```
px live -f service_endpoints_summary -- -start_time '-30m' -service 'default/echo-service'
```

This PxL sript takes a `service` argument. Note that Pixie formats service names in the `<namespace>/<service>` format.

<img src=".readme_assets/service_endpoints_summary.png" alt="Overview of endpoints for a service.">

To see timeseries graphs for endpoint latency, error and throughput, run the following command and then click the `Live View` link at the top:

```
px live pxbeta/service_endpoints -- -start_time '-30m' -service 'default/echo-service'
```

## Full-body HTTP/2 Requests for a Specified Service

From the top-level `endpoint-deprecation` folder, run:

```
px live -f service_requests -- -start_time '-30m' -service 'default/echo-service'
```

## Full-body HTTP/2 Requests for a Specified Logical Endpoint

From the top-level `endpoint-deprecation` folder, run:

```
px live -f service_endpoint_requests -- -start_time '-30m' -service 'default/echo-service' -endpoint '/v1/catalog/*/details'
```

<img src=".readme_assets/service_endpoint_requests.png" alt="Sample of requests sent to an endpoint.">

To inspect truncated table cells (e.g. `req_headers`), select the cell then press `enter`. To exit the expanded view, use `esc`.

## A List of Unique Request Header Field Values

From the top-level `endpoint-deprecation` folder, run:

```
px live -f unique_req_header_values -- -start_time '-30m' -service 'default/echo-service' -endpoint '/v1/catalog/*/details'
```

<img src=".readme_assets/unique_req_header_values.png" alt="List of unique request header field values.">

This script only examines two request header fields (Referer and Api-Key), but the `.pxl` file can be easily modified to inspect other fields.

## How to Interact with the Live CLI

For more information, check out the [reference docs](https://docs.px.dev/using-pixie/using-cli/#use-the-live-cli).

## How to Run Scripts in the Live UI

To learn how to run these scripts using the Live UI (instead of the CLI), check out the [reference docs](https://docs.px.dev/using-pixie/using-live-ui/#use-the-scratch-pad).

## Bugs & Features

Feel free to file a bug or an issue for a feature request. You can also join our [Slack](https://slackin.px.dev/) community.
