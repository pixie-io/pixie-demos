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
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

func launchTimeAfter() {
	go func() {
		<-time.After(10 * time.Second)
		fmt.Println("Waited for 10 seconds")
	}()
}

func launchSleep() {
	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println("Slept for 10 seconds")
	}()
}

// generateRandomString function adapted from
// https://kpbird.medium.com/golang-generate-fixed-size-random-string-dd6dbd5e63c0
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func generateRandomString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// Allocate a bunch of memory
func generateRandomStringArray(arrLength, strBytes int) []string {
	arr := make([]string, arrLength)
	for i := 0; i < arrLength; i++ {
		arr[i] = generateRandomString(strBytes)
	}
	return arr
}

func parseArrayArgs(r *http.Request) (int, int) {
	arrayLengthStr := r.URL.Query().Get("arrayLength")
	bytesPerElementStr := r.URL.Query().Get("bytesPerElement")
	arrayLength, err := strconv.Atoi(arrayLengthStr)
	if err != nil {
		arrayLength = 1000
	}
	bytesPerElement, err := strconv.Atoi(bytesPerElementStr)
	if err != nil {
		bytesPerElement = 8192
	}
	return arrayLength, bytesPerElement
}

func main() {
	http.HandleFunc("/allocate-memory-and-run-gc", func(w http.ResponseWriter, r *http.Request) {
		arrayLength, bytesPerElement := parseArrayArgs(r)
		arr := generateRandomStringArray(arrayLength, bytesPerElement)
		fmt.Fprintf(w, fmt.Sprintf("Generated string array with %d bytes of data\n", len(arr)*len(arr[0])))
		runtime.GC()
		fmt.Fprintf(w, "Ran garbage collector\n")
	})

	http.HandleFunc("/allocate-memory", func(w http.ResponseWriter, r *http.Request) {
		arrayLength, bytesPerElement := parseArrayArgs(r)
		arr := generateRandomStringArray(arrayLength, bytesPerElement)
		fmt.Fprintf(w, fmt.Sprintf("Generated string array with %d bytes of data\n", len(arr)*len(arr[0])))
	})

	http.HandleFunc("/run-gc", func(w http.ResponseWriter, r *http.Request) {
		runtime.GC()
		fmt.Fprintf(w, "Ran garbage collector\n")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
