Requirements:
  - You need a k8s cluster capable of running Pixie (see [here](https://docs.px.dev/installing-pixie/)).
  - Your local machine needs to have `envsubst` and `kubectl` installed with access to that k8s cluster.
  - If you want to also deploy a working stripe exfiltration example, you will need a Stripe account and a test api key for that account (api keys starting with `sk_test`)

Deployment Steps:
- First choose a public endpoint outside of your cluster to fake data exfiltration to. This can be any HTTP endpoint that will accept POST requests. Once you have that endpoint run:
  ```
  export EGRESS_URL=my-domain.tld/req/path
  ```
  replacing `my-domain.tld` with the domain and `/req/path` with the path to the endpoint you're going to be using.
- Now, if you have one, export your stripe test api key as an environment variable:
  ```
  export STRIPE_TEST_API_KEY=<your-api-key-here>
  ```
- [Deploy](https://docs.px.dev/installing-pixie/install-guides/community-cloud-for-pixie#4.-deploy-pixie) Pixie to your cluster.
- Deploy the demo by running the following in the same terminal you ran the `export` commands:
  ```
  kubectl create namespace px-data-exfiltration-demo
  envsubst < demo.yaml | kubectl apply -f -
  ```
- Explore the exfiltrated demo data in Pixie using the `px/cluster_egress` script.
