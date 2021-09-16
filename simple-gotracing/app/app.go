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
	"fmt"
	"net/http"
	"strconv"
)

// computeE computes the approximation of e by running a fixed number of iterations.
//go:noinline
func computeE(iterations int64) float64 {
	res := 2.0
	fact := 1.0

	for i := int64(2); i < iterations; i++ {
		fact *= float64(i)
		res += 1 / fact
	}
	return res
}

func main() {
	addr := ":9090"
	http.HandleFunc("/e", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		iters := int64(100)
		keys, ok := r.URL.Query()["iters"]
		if ok && len(keys[0]) >= 1 {
			val, err := strconv.ParseInt(keys[0], 10, 64)
			if err != nil || val <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			iters = val
		}

		w.Write([]byte(fmt.Sprintf("e = %0.4f\n", computeE(iters))))
	})

	fmt.Printf("Starting server on: %+v\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("Failed to run http server: %v\n", err)
	}
}
