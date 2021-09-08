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

#include "probe_deployment.h"

#include <bcc/BPF.h>
#include <linux/perf_event.h>

#include <memory>
#include <string>

int BCCWrapper::Init(const std::string& bpf_code, const std::vector<std::string>& cflags) {
  auto init_res = init(bpf_code, cflags);
  if (init_res.code() != 0) {
    std::cerr << "Unable to initialize BCC BPF program: " << init_res.msg() << std::endl;
    return 1;
  }
  return 0;
}

int BCCWrapper::AttachUProbe(const UProbeSpec& probe) {
  ebpf::StatusTuple attach_status = attach_uprobe(probe.obj_path, probe.symbol, probe.probe_fn,
                                                  /* address */ 0, probe.attach_type);
  if (attach_status.code() != 0) {
    std::cerr << "Failed to attach uprobe to binary " << probe.obj_path << " at symbol "
              << probe.symbol << ", error message: " << attach_status.msg() << std::endl;
    return 1;
  }
  std::cout << "Attached uprobe to binary " << probe.obj_path << " at symbol " << probe.symbol
            << std::endl;
  return 0;
}

int BCCWrapper::DetachUProbe(const UProbeSpec& probe) {
  ebpf::StatusTuple detach_status =
      detach_uprobe(probe.obj_path, probe.symbol, 0, probe.attach_type);

  if (detach_status.code() != 0) {
    std::cerr << "Failed to detach uprobe from binary " << probe.obj_path << " on symbol "
              << probe.symbol << ", error message: " << detach_status.msg() << std::endl;
    return 1;
  }
  return 0;
}

int BCCWrapper::OpenPerfBuffer(const PerfBufferSpec& perf_buffer, void* cb_cookie) {
  ebpf::StatusTuple open_status = open_perf_buffer(perf_buffer.name, perf_buffer.probe_output_fn,
                                                   perf_buffer.probe_loss_fn, cb_cookie);
  if (open_status.code() != 0) {
    std::cerr << "Failed to open perf buffer: " << perf_buffer.name
              << ", error message: " << open_status.msg() << std::endl;
    return 1;
  }
  std::cout << "Opened perf buffer " << perf_buffer.name << std::endl;
  return 0;
}

int BCCWrapper::ClosePerfBuffer(const PerfBufferSpec& perf_buffer) {
  ebpf::StatusTuple close_status = close_perf_buffer(std::string(perf_buffer.name));
  if (close_status.code() != 0) {
    std::cerr << "Failed to close perf buffer: " << perf_buffer.name
              << ", error message: " << close_status.msg() << std::endl;
    return 1;
  }
  return 0;
}

void BCCWrapper::PollPerfBuffer(const std::string& perf_buffer_name, int timeout_ms) {
  auto perf_buffer = get_perf_buffer(perf_buffer_name);
  if (perf_buffer != nullptr) {
    perf_buffer->poll(timeout_ms);
  }
}
