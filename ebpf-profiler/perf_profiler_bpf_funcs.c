/*
 * Copyright © 2018- Pixie Labs Inc.
 * Copyright © 2020- New Relic, Inc.
 * All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of New Relic Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Pixie Labs Inc. and its suppliers and
 * may be covered by U.S. and Foreign Patents, patents in process,
 * and are protected by trade secret or copyright law. Dissemination
 * of this information or reproduction of this material is strictly
 * forbidden unless prior written permission is obtained from
 * New Relic, Inc.
 *
 * SPDX-License-Identifier: Proprietary
 */

// LINT_C_FILE: Do not remove this line. It ensures cpplint treats this as a C file.

#include <linux/bpf_perf_event.h>
#include <linux/ptrace.h>

// NOLINTNEXTLINE: build/include_subdir
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
