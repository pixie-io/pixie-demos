/*
 * Copyright 2018- The Pixie Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

#include <signal.h>

#include <fstream>
#include <iostream>
#include <streambuf>
#include <string>

#include "openssl_tracer_types.h"
#include "probe_deployment.h"

// A probe on entry of SSL_write
UProbeSpec kSSLWriteEntryProbeSpec{
    .obj_path = "/usr/lib/x86_64-linux-gnu/libssl.so.1.1",
    .symbol = "SSL_write",
    .attach_type = BPF_PROBE_ENTRY,
    .probe_fn = "probe_entry_SSL_write",
};

// A probe on return of SSL_write
UProbeSpec kSSLWriteRetProbeSpec{
    .obj_path = "/usr/lib/x86_64-linux-gnu/libssl.so.1.1",
    .symbol = "SSL_write",
    .attach_type = BPF_PROBE_RETURN,
    .probe_fn = "probe_ret_SSL_write",
};

// A probe on entry of SSL_read
UProbeSpec kSSLReadEntryProbeSpec{
    .obj_path = "/usr/lib/x86_64-linux-gnu/libssl.so.1.1",
    .symbol = "SSL_read",
    .attach_type = BPF_PROBE_ENTRY,
    .probe_fn = "probe_entry_SSL_read",
};

// A probe on return of SSL_read
UProbeSpec kSSLReadRetProbeSpec{
    .obj_path = "/usr/lib/x86_64-linux-gnu/libssl.so.1.1",
    .symbol = "SSL_read",
    .attach_type = BPF_PROBE_RETURN,
    .probe_fn = "probe_ret_SSL_read",
};

const std::vector<UProbeSpec> kUProbes = {
    kSSLWriteEntryProbeSpec,
    kSSLWriteRetProbeSpec,
    kSSLReadEntryProbeSpec,
    kSSLReadRetProbeSpec,
};

void handle_output(void* /*cb_cookie*/, void* data, int /*data_size*/) {
  // Copy the raw memory into the ssl_data_event_t struct that we know it is.
  // This also addresses any memory alignment issues, as the raw bytes are not necessarily aligned.
  struct ssl_data_event_t r;
  std::memcpy(&r, data, sizeof(r));

  std::string_view plaintext(r.data, r.data_len);

  std::cout << " t=" << r.timestamp_ns;
  std::cout << " type=" << (r.type == kSSLRead ? "read" : "write");
  std::cout << " data=" << plaintext;
  std::cout << std::endl;
}

const PerfBufferSpec kPerfBufferSpec = {
    .name = "tls_events",
    .probe_output_fn = &handle_output,
    .probe_loss_fn = nullptr,
};

#define RETURN_IF_ERROR(x) \
  if (x != 0) return 1;

int main(int argc, char** argv) {
  // Read arguments to get the target PID to trace.
  if (argc != 2) {
    std::cerr << "Usage: " << argv[0] << " <PID to trace for SSL traffic>" << std::endl;
    exit(1);
  }
  std::string target_pid(argv[1]);

  BCCWrapper bcc;

  // Read the BPF code.
  std::ifstream ifs("openssl_tracer_bpf_funcs.c");
  std::string bpf_code(std::istreambuf_iterator<char>(ifs), {});

  // Compile the BPF code and load into the kernel.
  // DTRACE_PID is a macro in the BPF code that controls the PID to be traced for SSL traffic.
  RETURN_IF_ERROR(bcc.Init(bpf_code, {"-DTRACE_PID=" + target_pid}));

  // Deploy uprobes.
  for (auto& probe_spec : kUProbes) {
    RETURN_IF_ERROR(bcc.AttachUProbe(probe_spec));
  }

  // Open the perf buffer used by our uprobes to output data.
  RETURN_IF_ERROR(bcc.OpenPerfBuffer(kPerfBufferSpec));

  std::cout << "Successfully deployed BPF probes. Tracing for SSL data. Use Ctrl-C to exit."
            << std::endl;

  // Periodically read the output buffer and print entries to screen.
  while (true) {
    bcc.PollPerfBuffer(kPerfBufferSpec.name);
    usleep(100000);
  }

  return 0;
}
