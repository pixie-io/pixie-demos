# OpenSSL Tracer using BPF

This is a basic example of how to trace the OpenSSL library using eBPF.
This tracer uses BCC to deploy the eBPF probes.
This demo was created to accompany the "Debugging with eBPF Part 3: Tracing SSL/TLS connections" [blog post](https://blog.px.dev/ebpf-openssl-tracing/).

## Prerequisites

You must have the BCC development package installed. On Ubuntu, the package can be installed as follows:

```
sudo apt install libbpfcc-dev
```

Other distributions have similar commands.

## Build

To compile, execute the following command:

```
make
```

## Run Demo Application

A demo application to trace is included. It is a simple client-server written in Python, which uses OpenSSL.

First, you'll have to generate some certificates for the client and server.
To keep things simple, you can generate some self-signed certificates as follows:

```
make -C ssl_client_server certs
```

To run the demo app, you'll need two terminals.

In one terminal, run the server:

```
cd ssl_client_server; ./server.py
```

In the second terminal, run the client:

```
cd ssl_client_server; ./client.py
```

## Run Tracer

The BPF tracer is run as follows:

```
sudo ./openssl_tracer <pid>
```

To run it on the demo app, run the following command in a separate terminal:

```
sudo ./openssl_tracer $(pgrep -f "./client.py")
```
