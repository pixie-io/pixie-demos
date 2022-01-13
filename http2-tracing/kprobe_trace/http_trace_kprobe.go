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

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"unsafe"

	"github.com/fatih/color"
	"github.com/iovisor/gobpf/bcc"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

var (
	tracePID     int
	parseHttp2   bool
	printEnabled bool
)

func init() {
	flag.IntVar(&tracePID, "pid", -1, "The pid to trace")
	flag.BoolVar(&parseHttp2, "parseHttp2", false, "If true, parse the data as HTTP2 frames")
	flag.BoolVar(&printEnabled, "print", true, "Print output")
}

// Track the event type. Must match consts in the C file.
type EventType int32

const (
	// Addr Event (accept).
	ETSyscallAddr EventType = iota + 1
	// Write event.
	ETSyscallWrite
	// Close Event.
	ETSyscallClose
)

// The attributes of the SyscallEvent. Must match consts in the C file.
type Attributes struct {
	EvType  EventType
	Fd      int32
	Bytes   int32
	MsgSize int32
}

// SyscallWriteEvent is the BPF struct for a syscall.
type SyscallWriteEvent struct {
	Attr Attributes
	Msg  []byte
}

// MessageInfo stores a buffer for partial messages on a specific file descriptor.
type MessageInfo struct {
	SocketInfo []byte
	Buf        bytes.Buffer
}

func mustAttachKprobeToSyscall(m *bcc.Module, probeType int, syscallName string, probeName string) {
	fnName := bcc.GetSyscallFnName(syscallName)
	kprobe, err := m.LoadKprobe(probeName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to attach probe: %s\n", err)
		os.Exit(1)
	}

	if probeType == bcc.BPF_PROBE_ENTRY {
		err = m.AttachKprobe(fnName, kprobe, -1)
	} else {
		err = m.AttachKretprobe(fnName, kprobe, -1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to attach entry probe %s: %s\n", probeName, err)
		os.Exit(1)
	}
}

type requestHandler struct {
	FdMap map[int32]*MessageInfo
}

// HandleBPFEvent handles an event from the BPF trace code.
func (r *requestHandler) HandleBPFEvent(v []byte) {
	var ev SyscallWriteEvent
	// To save space on the perf buffer we write out only the number of bytes necessary for the request.
	// The attributes are used to figure out how big the message is, so we start by reading them first.
	// After that we can allocate a buffer to read the data.
	if err := binary.Read(bytes.NewBuffer(v), bcc.GetHostByteOrder(), &ev.Attr); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode struct: %+v\n", err)
		return
	}
	// Now read the actual data.
	ev.Msg = make([]byte, ev.Attr.MsgSize)
	if err := binary.Read(bytes.NewBuffer(v[unsafe.Sizeof(ev.Attr):]), bcc.GetHostByteOrder(), &ev.Msg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode struct: %+v\n", err)
		return
	}

	// Based on the event type we:
	//   Insert on AcceptEvent.
	//   Append on Msg event.
	//   Delete and print on close.
	switch ev.Attr.EvType {
	case ETSyscallAddr:
		r.FdMap[ev.Attr.Fd] = &MessageInfo{
			SocketInfo: ev.Msg,
		}
	case ETSyscallWrite:
		if elem, ok := r.FdMap[ev.Attr.Fd]; ok {
			elem.Buf.Write(ev.Msg)
		}
	case ETSyscallClose:
		if msgInfo, ok := r.FdMap[ev.Attr.Fd]; ok {
			delete(r.FdMap, ev.Attr.Fd)

			// Decode in a go routine to avoid blocking the perf buffer read.
			go parseAndPrintMessage(msgInfo)
		} else {
			fmt.Fprintf(os.Stderr, "Missing request with FD: %d\n", ev.Attr.Fd)
			return
		}
	}
}

func parseHttpMessages(msgInfo *MessageInfo) {
	// We have the complete request so we try to parse the actual HTTP request.
	resp, err := http.ReadResponse(bufio.NewReader(&msgInfo.Buf), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse request, err: %v\n", err)
		return
	}

	if printEnabled {
		body := resp.Body
		b, _ := ioutil.ReadAll(body)
		body.Close()
		fmt.Printf("StatusCode: %s, Len: %s, ContentType: %s, Body: %s\n",
			color.GreenString("%d", resp.StatusCode),
			color.GreenString("%d", resp.ContentLength),
			color.GreenString("%s", resp.Header["Content-Type"]),
			color.GreenString("%s", string(b)))
	}
}

func formatFrame(f http2.Frame, decoder *hpack.Decoder) string {
	var buf bytes.Buffer
	switch f := f.(type) {
	case *http2.DataFrame:
		fmt.Fprintf(&buf, "[DATA] %q", f.Data())
	case *http2.HeadersFrame:
		fmt.Fprintf(&buf, "[HEADERS]")
		str := hex.EncodeToString(f.HeaderBlockFragment())
		log.Println("Hex header code: " + str)
		hfs, _ := decoder.DecodeFull(f.HeaderBlockFragment())
		for _, hf := range hfs {
			fmt.Fprintf(&buf, " %q:%q", hf.Name, hf.Value)
		}
	default:
		fmt.Fprintf(&buf, "[IGNORED] %v", f)
	}
	return buf.String()
}

func parseHttp2Frames(msgInfo *MessageInfo) {
	framer := http2.NewFramer(ioutil.Discard, &msgInfo.Buf)
	decoder := hpack.NewDecoder(2048, nil)
	for {
		f, err := framer.ReadFrame()
		if err != nil {
			log.Println("ReadFrame failed, err=%v", err)
			break
		}
		log.Println(formatFrame(f, decoder))
	}
}

func parseAndPrintMessage(msgInfo *MessageInfo) {
	if parseHttp2 {
		parseHttp2Frames(msgInfo)
		return
	}
	parseHttpMessages(msgInfo)
}

func main() {
	flag.Parse()
	if tracePID < 0 {
		panic("Argument --pid needs to be specified")
	}
	bpfProgramResolved := strings.ReplaceAll(bpfProgram, "$PID", fmt.Sprintf("%d", tracePID))
	bccMod := bcc.NewModule(bpfProgramResolved, []string{})
	mustAttachKprobeToSyscall(bccMod, bcc.BPF_PROBE_ENTRY, "accept4", "syscall__probe_entry_accept4")
	mustAttachKprobeToSyscall(bccMod, bcc.BPF_PROBE_RETURN, "accept4", "syscall__probe_ret_accept4")
	mustAttachKprobeToSyscall(bccMod, bcc.BPF_PROBE_ENTRY, "write", "syscall__probe_write")
	mustAttachKprobeToSyscall(bccMod, bcc.BPF_PROBE_ENTRY, "close", "syscall__probe_close")

	// Create the output table named "golang_http_response_events" that the BPF program writes to.
	table := bcc.NewTable(bccMod.TableId("syscall_write_events"), bccMod)
	ch := make(chan []byte)

	pm, err := bcc.InitPerfMap(table, ch, nil)
	if err != nil {
		panic(err)
	}

	// Watch Ctrl-C so we can quit this program.
	intCh := make(chan os.Signal, 1)
	signal.Notify(intCh, os.Interrupt)

	pm.Start()
	defer pm.Stop()

	// Map from file descriptor the MessageInfo.
	requestHander := &requestHandler{
		FdMap: make(map[int32]*MessageInfo, 0),
	}

	for {
		select {
		case <-intCh:
			fmt.Println("Terminating")
			os.Exit(0)
		case v := <-ch:
			requestHander.HandleBPFEvent(v)
		}
	}
}
