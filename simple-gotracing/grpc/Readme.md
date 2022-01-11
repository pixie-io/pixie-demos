# gRPC client & server

This directory contains gRPC client & server for demonstrating the kprobe- & uprobe-based HTTP2
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
(cd .. && go build -o /tmp/http_trace_kprobe ./http_trace_kprobe && \
 sudo -E /tmp/http_trace_kprobe --pid=${PID} --parseHttp2)

# Build and launch the uprobe tracer.
(cd .. && go build -o /tmp/http2_trace_uprobe ./http2_trace_uprobe && \
 sudo -E /tmp/http2_trace_uprobe --binary=/tmp/grpc_server)

# Build and launch gRPC client.
go build -o /tmp/grpc_client client/main.go && /tmp/grpc_client --count 10
```
