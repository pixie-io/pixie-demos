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
* [OPTIONAL] Change `gobpf` version in go.mod if `go build` failed building uprobe tracer with error
  like below:
  ```
  # github.com/iovisor/gobpf/bcc
  ../../../pkg/mod/github.com/iovisor/gobpf@v0.0.0-20200614202714-e6b321d32103/bcc/module.go:261:33: not enough arguments in call to _C2func_bpf_attach_uprobe
        have (_Ctype_int, uint32, *_Ctype_char, *_Ctype_char, _Ctype_ulong, _Ctype_int)
        want (_Ctype_int, uint32, *_Ctype_char, *_Ctype_char, _Ctype_ulong, _Ctype_int, _Ctype_uint)
  ```
  You can use the gobpf version in the comments to replace the gobpf version.

## Usage

```
# To update protobuf generated go source files:
protoc --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false proto/greet.proto

# Build and launch gRPC server.
go build -o /tmp/grpc_server server/main.go && /tmp/grpc_server

# Get the PID of the gRPC server.
PID=$(ps aux | grep grpc_server | head -n 1 | awk '{print $2}')

# Build and launch the uprobe tracer.
go build -o /tmp/uprobe_trace ./uprobe_trace && \
 sudo -E /tmp/uprobe_trace --binary=/tmp/grpc_server

# Build and launch gRPC client.
go build -o /tmp/grpc_client client/main.go && /tmp/grpc_client --count 10
```

## Bugs & Features

Feel free to file a bug or an issue for a feature request. You can also join our
[Slack](https://slackin.px.dev/) community.
