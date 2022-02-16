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

#include <linux/sched.h>

#include "types.h"

// Effectively returns task->group_leader->real_start_time;
// Note that before Linux 5.5, real_start_time was called start_boottime.
static inline __attribute__((__always_inline__)) uint64_t get_tgid_start_time() {
  struct task_struct* task = (struct task_struct*)bpf_get_current_task();
  struct task_struct* group_leader_ptr = task->group_leader;
  uint64_t start_time = group_leader_ptr->start_time;
  return div_u64(start_time, NSEC_PER_SEC / USER_HZ);
}

// Map that stores counts of recv calls by PID
BPF_HASH(counts_by_pid, struct tgid_ts_t, int64_t);

// Probe that counts every time it is triggered.
// Can be used to count things like syscalls or particular functions.
int syscall__probe_counter(struct pt_regs* ctx) {
  uint32_t tgid = bpf_get_current_pid_tgid() >> 32;
  struct tgid_ts_t process_id = {};
  process_id.tgid = tgid;
  process_id.ts = get_tgid_start_time();

  int64_t kInitVal = 0;
  int64_t* count = counts_by_pid.lookup_or_init(&process_id, &kInitVal);
  if (count != NULL) {
    *count = *count + 1;
  }

  return 0;
}

