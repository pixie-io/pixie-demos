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

#pragma once

#include <bcc/BPF.h>
#include <linux/perf_event.h>

#include <iostream>
#include <map>
#include <string>
#include <vector>

/**
 * Describes a userspace probe (uprobe).
 */
struct UProbeSpec {
  // The path to the object file (e.g. binary, shared object) to which this uprobe is attached.
  std::string obj_path;

  // Symbol within the binary to probe.
  std::string symbol;

  // Whether to attach on entry or return of function.
  bpf_probe_attach_type attach_type;

  // BPF function to execute when the uprobe triggers.
  std::string probe_fn;
};

/**
 * Describes a BPF perf buffer, through which data is returned to user-space.
 */
struct PerfBufferSpec {
  // Name of the perf buffer.
  // Must be the same as the perf buffer name declared in the probe code with BPF_PERF_OUTPUT.
  std::string name;

  // Function that will be called for every event in the perf buffer,
  // when perf buffer read is triggered.
  perf_reader_raw_cb probe_output_fn;

  // Function that will be called if there are lost/clobbered perf events.
  perf_reader_lost_cb probe_loss_fn;
};

class BCCWrapper : public ebpf::BPF {
 public:
  int Init(const std::string& bpf_code, const std::vector<std::string>& cflags);
  int AttachUProbe(const UProbeSpec& probe);
  int DetachUProbe(const UProbeSpec& probe);
  int OpenPerfBuffer(const PerfBufferSpec& perf_buffer, void* cb_cookie = nullptr);
  int ClosePerfBuffer(const PerfBufferSpec& perf_buffer);
  void PollPerfBuffer(const std::string& perf_buffer_name, int timeout_ms = 1);
};
