# Copyright © 2018- Pixie Labs Inc.
# Copyright © 2020- New Relic, Inc.
# All Rights Reserved.
#
# NOTICE:  All information contained herein is, and remains
# the property of New Relic Inc. and its suppliers,
# if any.  The intellectual and technical concepts contained
# herein are proprietary to Pixie Labs Inc. and its suppliers and
# may be covered by U.S. and Foreign Patents, patents in process,
# and are protected by trade secret or copyright law. Dissemination
# of this information or reproduction of this material is strictly
# forbidden unless prior written permission is obtained from
# New Relic, Inc.
#
# SPDX-License-Identifier: Proprietary

all: perf_profiler sqrt

perf_profiler: perf_profiler.cc
	clang++ --std=c++17 -o $@ $^ -lbcc

sqrt: toy_app/sqrt.go
	go build $^

clean:
	rm perf_profiler sqrt
