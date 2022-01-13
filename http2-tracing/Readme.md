# gRPC client & server

This directory contains gRPC client & server for demonstrating the kprobe & uprobe-based HTTP2
tracers. The following shell commands have to be run from this directory.

```
# You might need to install go protobuf plugin: https://grpc.io/docs/languages/go/quickstart/
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
