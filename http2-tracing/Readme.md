# HTTP/2 Tracing With eBPF Uprobes Demo

Use eBPF uprobes to trace HTTP/2 headers, without any changes to the application code.

## What is this demo?

This demo provides the gRPC client and server, and the uprobe tracer for
[HTTP2 tracing](https://blog.px.dev/http2-tracing).

## Prerequisites

* This demo only works on Linux, and with eBPF support. The code was tested on Ubuntu 20.04.3 LTS
  with 5.4 kernel.
* Install [BCC](https://github.com/iovisor/bcc/blob/master/INSTALL.md).
* Install [Protocol buffer compiler](https://grpc.io/docs/protoc-installation/) and
  [go protobuf plugin](https://grpc.io/docs/languages/go/quickstart/).
* Various Go packages might be installed, follow the directives when any of the `go build` commands
  failed.

## Usage

```
# To update protobuf generated go source files:
protoc --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false proto/greet.proto

# Build and launch gRPC server.
go build -o /tmp/grpc_server server/main.go && /tmp/grpc_server

# Get the PID of the gRPC server.
PID=$(ps aux | grep grpc_server | head -n 1 | awk '{print $2}')

# Build and launch the kprobe tracer.
go build -o /tmp/kprobe_trace ./kprobe_trace && \
 sudo -E /tmp/kprobe_trace --pid=${PID}

# Build and launch the uprobe tracer.
go build -o /tmp/uprobe_trace ./uprobe_trace && \
 sudo -E /tmp/uprobe_trace --binary=/tmp/grpc_server

# Build and launch gRPC client.
go build -o /tmp/grpc_client client/main.go && /tmp/grpc_client --count 10
```
