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

package main

// This does not work for Golang programs built with 1.17 or newer toolchain. Go 1.17 uses a register-based calling
// convention which the BPF code here cannot handle.
const bpfProgram = `
#include <uapi/linux/ptrace.h>

#define HEADER_FIELD_STR_SIZE 128
#define MAX_HEADER_COUNT 64

BPF_PERF_OUTPUT(go_http2_header_events);

struct header_field_t {
  int32_t size;
  char msg[HEADER_FIELD_STR_SIZE];
};

struct go_grpc_http2_header_event_t {
  struct header_field_t name;
  struct header_field_t value;
};

// This matches the golang string object memory layout. Used to help read golang string objects in BPF code.
struct gostring {
  const char* ptr;
  int64_t len;
};

static int64_t min(int64_t l, int64_t r) {
  return l < r ? l : r;
}

// Copy the content of a hpack.HeaderField object into header_field_t object.
static void copy_header_field(struct header_field_t* dst, const void* header_field_ptr) {
  struct gostring str = {};
  bpf_probe_read(&str, sizeof(str), header_field_ptr);
  if (str.len <= 0) {
    dst->size = 0;
    return;
  }
  dst->size = min(str.len, HEADER_FIELD_STR_SIZE);
  bpf_probe_read(dst->msg, dst->size, str.ptr);
}

// Copies and submits content of an array of hpack.HeaderField to perf buffer.
static void submit_headers(struct pt_regs* ctx, void* fields_ptr, int64_t fields_len) {
  // Size of the golang hpack.HeaderField struct.
  const size_t header_field_size = 40;
  struct go_grpc_http2_header_event_t event = {};
  for (size_t i = 0; i < MAX_HEADER_COUNT; ++i) {
    if (i >= fields_len) {
      continue;
    }
    const void* header_field_ptr = fields_ptr + i * header_field_size;
    copy_header_field(&event.name, header_field_ptr);
    copy_header_field(&event.value, header_field_ptr + 16);
    go_http2_header_events.perf_submit(ctx, &event, sizeof(event));
  }
}

// Signature: func (l *loopyWriter) writeHeader(streamID uint32, endStream bool, hf []hpack.HeaderField, onWrite func())
int probe_loopy_writer_write_header(struct pt_regs* ctx) {
  const void* sp = (const void*)ctx->sp;

  void* fields_ptr;
	const int kFieldsPtrOffset = 24;
  bpf_probe_read(&fields_ptr, sizeof(void*), sp + kFieldsPtrOffset);

  int64_t fields_len;
	const int kFieldsLenOffset = 8;
  bpf_probe_read(&fields_len, sizeof(int64_t), sp + kFieldsPtrOffset + kFieldsLenOffset);

  submit_headers(ctx, fields_ptr, fields_len);
  return 0;
}

// Signature: func (t *http2Server) operateHeaders(frame *http2.MetaHeadersFrame, handle func(*Stream),
// traceCtx func(context.Context, string) context.Context)
int probe_http2_server_operate_headers(struct pt_regs* ctx) {
  const void* sp = (const void*)ctx->sp;

  void* frame_ptr;
  bpf_probe_read(&frame_ptr, sizeof(void*), sp + 16);

  void* fields_ptr;
  bpf_probe_read(&fields_ptr, sizeof(void*), frame_ptr + 8);

  int64_t fields_len;
  bpf_probe_read(&fields_len, sizeof(int64_t), frame_ptr + 8 + 8);

  submit_headers(ctx, fields_ptr, fields_len);
  return 0;
}
`
