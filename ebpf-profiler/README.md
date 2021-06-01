
# Performance Profiler using eBPF

This is a demo performance profiler, written to accompany the Pixie [blog post](https://blog.px.dev/cpu-profiling).
It shows how to get sample stack traces for performance profiling, using eBPF.

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

A demo application to profile is included: `toy_app/sqrt.go`

## Run Tracer

The BPF performance profiler is run as follows:

```
sudo ./perf_profiler <target pid> <duration>
```

To run it on the demo app for 30 seconds, run the following command in a separate terminal:
```
sudo ./perf_profiler $(pgrep -f "sqrt") 30
```

