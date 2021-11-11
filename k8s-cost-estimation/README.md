# Kubernetes Cluster Cost Estimation

These PxL scripts demonstrate how to use telemetry data automatically provided by Pixie to estimate the cost of your Kubernetes cluster.

- **cpu_cost.pxl**: Estimates yearly CPU cost per Service based on last hour of usage.
- **mem_cost.pxl**: Estimates yearly memory cost per Service based on last hour of usage.
- **network_cost.pxl**: Estimates yearly network cost per Service based on last hour of usage.
- **cost_per_request.pxl**: Estimates number of requests per Service based on last hour of usage. Divides total cost (sum of CPU, memory, network costs) by number of requests per Service.

## Prerequisites

You will need a Kubernetes cluster with Pixie installed. If you do not have a cluster, you can create a minikube cluster and install Pixie using one of our [install guides](https://docs.px.dev/installing-pixie/install-guides/).

## Instructions

1. Adjust the cost placeholder values at the top of the scripts depending on your cloud provider costs.

2. Run the scripts using Pixie's Live UI or CLI.

Using the Live UI:

- Select the `Scratch Pad` script from the script drop-down menu.
- Open the editor using `ctrl+e` (Windows, Linux) or `cmd+e`(Mac).
- Paste the contents of the script into the editor and close the editor using `ctrl+e` (Windows, Linux) or `cmd+e`(Mac).
- Press the `RUN` button (top right) to execute the script.
- Sort table columns by clicking the column title.

Using the CLI:

- Copy the scripts to a local directory.
- Run `px run -f <script.pxl>`

## Have questions? Need help?

Please reach out on our Pixie Community [Slack](https://slackin.px.dev/) or file a GitHub issue.
