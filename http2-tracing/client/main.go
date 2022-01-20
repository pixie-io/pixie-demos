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
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "pixielabs.ai/pixie/demos/http2-tracing/proto/greetpb"
)

func mustCreateGrpcClientConn(address string) *grpc.ClientConn {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

func connectAndGreet(address, name string, count, sleep_millis int) {
	conn := mustCreateGrpcClientConn(address)
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	for i := 0; i < count; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.Message)
		time.Sleep(time.Duration(sleep_millis) * time.Millisecond)
	}
}

func main() {
	address := flag.String("address", "localhost:50051", "Server end point.")
	name := flag.String("name", "world", "The name to greet.")
	count := flag.Int("count", 1, "The number of RPC calls to make.")
	sleep_millis := flag.Int("sleep-millis", 500, "The number of milliseconds to sleep between RPC calls.")
	flag.Parse()
	connectAndGreet(*address, *name, *count, *sleep_millis)
}
