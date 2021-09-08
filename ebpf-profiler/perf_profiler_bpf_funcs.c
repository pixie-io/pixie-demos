/*
 * Copyright 2018- The Pixie Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * 
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include <linux/bpf_perf_event.h>
#include <linux/ptrace.h>

#include "perf_profiler_types.h"

const int kNumMapEntries = 65536;

BPF_STACK_TRACE(stack_traces, kNumMapEntries);

BPF_HASH(histogram, struct stack_trace_key_t, uint64_t, kNumMapEntries);

int sample_stack_trace(struct bpf_perf_event_data* ctx) {
  // Sample the user stack trace, and record in the stack_traces structure.
  int user_stack_id = stack_traces.get_stackid(&ctx->regs, BPF_F_USER_STACK);

  // Sample the kernel stack trace, and record in the stack_traces structure.
  int kernel_stack_id = stack_traces.get_stackid(&ctx->regs, 0);

  // Update the counters for this user+kernel stack trace pair.
  struct stack_trace_key_t key = {};
  key.pid = bpf_get_current_pid_tgid() >> 32;
  key.user_stack_id = user_stack_id;
  key.kernel_stack_id = kernel_stack_id;
  histogram.increment(key);

  return 0;
}
