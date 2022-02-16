# eBPF Portability Demos

This folder contains a basic eBPF probes to demonstrate portability issues.

There are two folders in this demo:
 - `basic`: A simple probe that counts the syscall of your choice, broken down per PID. This simple probe has no dependencies on Linux headers and is generally robust.
 - `not_portable`: A variation of the basic probe that also considers the PID start time. This probe includes the Linux sched.h and would therefore be brittle if compiled with the wrong Linux headers.

