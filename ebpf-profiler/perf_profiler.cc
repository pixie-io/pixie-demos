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

#include <bcc/BPF.h>
#include <linux/perf_event.h>

#include <fstream>
#include <iostream>
#include <string>

#include "perf_profiler_types.h"

//-----------------------------------------------------------------------------
// Global constants/parameters.
//-----------------------------------------------------------------------------

// The name of the file that contains our BPF code.
constexpr char kBPFFileName[] = "perf_profiler_bpf_funcs.c";

// The name of the function in our BPF code that collects a stack trace when triggered.
constexpr char kProbeFn[] = "sample_stack_trace";

// The period with which we want to collect stack traces.
constexpr uint64_t kSamplingPeriodMillis = 11;

// The name of the stack traces map inside our BPF code.
// This is where stack traces are stored.
constexpr char kStackTracesMapName[] = "stack_traces";

// The name of the histogram map inside our BPF code.
// This maps stack traces to the number of times they were observed.
constexpr char kHistogramMapName[] = "histogram";

//-----------------------------------------------------------------------------

#define RETURN_IF_ERROR(x) \
  if (x != 0) return 1;

/**
 * Loads the provided BPF program into the kernel.
 */
int InitBPFProgram(ebpf::BPF* bcc, const std::string& bpf_code) {
  auto init_res = bcc->init(bpf_code);
  if (init_res.code() != 0) {
    std::cerr << "Unable to initialize BCC BPF program: " << init_res.msg() << std::endl;
    return 1;
  }
  return 0;
}

/**
 * Set up a periodic event to regularly trigger a BPF function.
 */
int AttachSamplingProbe(ebpf::BPF* bcc, std::string_view probe_fn,
                        uint64_t sampling_period_millis) {
  constexpr uint64_t kNanosPerMilli = 1000 * 1000;

  // A sampling probe is just a perf event probe, where the perf event is a clock counter.
  // When a requisite number of clock samples occur, the kernel will trigger the BPF code.
  // By specifying a frequency, the kernel will attempt to adjust the threshold to achieve
  // the desired sampling frequency.
  ebpf::StatusTuple attach_status =
      bcc->attach_perf_event(PERF_TYPE_SOFTWARE, PERF_COUNT_SW_CPU_CLOCK, std::string(probe_fn),
                             sampling_period_millis * kNanosPerMilli, 0);
  if (attach_status.code() != 0) {
    std::cerr << "Failed to attach perf_event: " << attach_status.msg() << std::endl;
    return 1;
  }
  return 0;
}

/**
 * Aggregate the collected stack trace samples into a map of stack traces and counts.
 * The stack traces retrieved from BPF are essentially vectors of addresses,
 * so this function turns those addresses into function symbols when possible.
 * The stack trace format is semi-colon delimited.
 */
std::map<std::string, int> ProcessStackTraces(ebpf::BPF* bcc, int target_pid) {
  ebpf::BPFStackTable stack_traces = bcc->get_stack_table(kStackTracesMapName);
  ebpf::BPFHashTable<struct stack_trace_key_t, uint64_t> histogram =
      bcc->get_hash_table<struct stack_trace_key_t, uint64_t>(kHistogramMapName);

  std::map<std::string, int> result;

  for (const auto& [key, count] : histogram.get_table_offline()) {
    if (key.pid != target_pid) {
      continue;
    }

    std::string stack_trace_str;

    if (key.user_stack_id >= 0) {
      std::vector<std::string> user_stack_symbols =
          stack_traces.get_stack_symbol(key.user_stack_id, key.pid);
      for (const auto& sym : user_stack_symbols) {
        stack_trace_str += sym;
        stack_trace_str += ";";
      }
    }

    if (key.kernel_stack_id >= 0) {
      std::vector<std::string> user_stack_symbols =
          stack_traces.get_stack_symbol(key.kernel_stack_id, -1);
      for (const auto& sym : user_stack_symbols) {
        stack_trace_str += sym;
        stack_trace_str += ";";
      }
    }

    result[stack_trace_str] += 1;
  }

  return result;
}

void PrintResults(const std::map<std::string, int>& stack_traces) {
  for (const auto& [stack_trace, count] : stack_traces) {
    std::cout << count << " " << stack_trace << std::endl;
  }
}

int main(int argc, char** argv) {
  // Read arguments to get the target PID to trace.
  if (argc != 3) {
    std::cerr << "Usage: " << argv[0] << " <PID to profile> <duration in seconds>" << std::endl;
    exit(1);
  }
  int target_pid = std::atoi(argv[1]);
  int duration_secs = std::atoi(argv[2]);

  ebpf::BPF bcc;

  // Read the BPF code.
  std::ifstream ifs(kBPFFileName);
  std::string bpf_code(std::istreambuf_iterator<char>(ifs), {});

  // Compile the BPF code and load into the kernel.
  RETURN_IF_ERROR(InitBPFProgram(&bcc, bpf_code));

  // Deploy stack trace sampling probe.
  RETURN_IF_ERROR(AttachSamplingProbe(&bcc, kProbeFn, kSamplingPeriodMillis));

  std::cout << "Successfully deployed BPF profiler." << std::endl;
  std::cout << "Collecting stack trace samples for " << duration_secs << " seconds." << std::endl;

  // Sleep for some time to allow stack traces to be sampled.
  sleep(duration_secs);

  // Now process and print the results.
  std::map<std::string, int> stack_traces = ProcessStackTraces(&bcc, target_pid);
  PrintResults(stack_traces);

  return 0;
}
