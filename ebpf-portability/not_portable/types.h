#pragma once

struct tgid_ts_t {
  uint32_t tgid;
  uint64_t ts; // Timestamp when the process started.
};
