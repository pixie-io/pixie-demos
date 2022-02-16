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

#include <fstream>
#include <iostream>
#include <streambuf>
#include <string>

#include <bcc/BPF.h>

// Counts number of syscalls made by each PID for the specified run time.
int main(int argc, char** argv) {
  if (argc != 3) {
    std::cout << "Usage: " << argv[0] << " <runtime> <syscall to count>" << std::endl;
    std::cout << "Example: " << argv[0] << " 5 recvmsg" << std::endl;
    exit(1);
  }

  int runtime = std::atoi(argv[1]);
  std::string syscall(argv[2]);

  // Read the BPF code.
  std::ifstream ifs("probes.c");
  std::string bpf_code(std::istreambuf_iterator<char>(ifs), {});

  ebpf::BPF bcc;

  bcc.init(std::string(bpf_code));

  auto fnname = bcc.get_syscall_fnname(syscall);
  bcc.attach_kprobe(fnname, "syscall__probe_counter");

  sleep(runtime);

  auto counts_by_pid = bcc.get_hash_table<uint32_t, int64_t>("counts_by_pid");

  std::cout << "PID: count" << std::endl;
  for (const auto& [pid, count] : counts_by_pid.get_table_offline()) {
    std::cout << pid << ":\t" << count << std::endl;
  }

  return 0;
}
