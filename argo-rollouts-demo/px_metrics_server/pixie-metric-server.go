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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/julienschmidt/httprouter"
	"px.dev/pxapi"
	pxTypes "px.dev/pxapi/types"
)

// PxL script to compute the metrics. Could be extended to compute additional metrics.
const timeWindow = "-30s"
const pxlScript = `import px

POD_NAMESPACE="%s"
POD_NAME="%s"
START_TIME="%s"

# Get HTTP events (not all pods will have this)
df = px.DataFrame(table='http_events', start_time=START_TIME)

# Add context
df.namespace = df.ctx['namespace']
df.pod = df.ctx['pod']

# Filter HTTP events
df = df[df.trace_role == 2]
df.failure = df.resp_status >= 400
df = df[df.namespace == POD_NAMESPACE]
df = df[px.contains(df.pod, POD_NAME)]

# Aggregate throughput, errors, latency for inbound requests to matching pods.
df = df.agg(
    http_req_count_in=('latency', px.count),
    http_error_count_in=('failure', px.sum),
    http_latency_in=('latency', px.quantiles)
)

df.pod=POD_NAME
df.latency_p99 = px.DurationNanos(px.floor(px.pluck_float64(df.http_latency_in, 'p99')))
df.http_error_rate_in = px.Percent(
        px.select(df.http_req_count_in != 0, df.http_error_count_in / df.http_req_count_in, 0.0))

px.display(df[['pod', 'http_error_rate_in']], 'pod_stats')
`

type pixieMetricsProvider struct {
	vizierClient *pxapi.VizierClient
	dataMux      sync.Mutex
	podErrorRate map[string]float64
}

func newPixieMetricProvider(apiKey string, cloudAddr string, clusterID string) *pixieMetricsProvider {
	// Create Pixie client.
	log.Println("Creating Pixie client.")
	ctx := context.Background()
	opts := []pxapi.ClientOption{pxapi.WithAPIKey(apiKey)}
	if cloudAddr != "" {
		opts = append(opts, pxapi.WithCloudAddr(cloudAddr))
	}
	pixieClient, err := pxapi.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalln(err.Error())
	}
	vz, err := pixieClient.NewVizierClient(ctx, clusterID)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Create pixieMetricsProvider.
	provider := &pixieMetricsProvider{
		vizierClient: vz,
		podErrorRate: make(map[string]float64),
	}

	return provider
}

func (p *pixieMetricsProvider) computeMetrics(ctx context.Context, pxlScript string) {
	tm := &tableMux{
		onPodStatsComplete: func(newStats map[string]float64) {
			p.dataMux.Lock()
			defer p.dataMux.Unlock()
			p.podErrorRate = newStats
		},
	}
	log.Println("Executing PxL query.")
	results, err := p.vizierClient.ExecuteScript(ctx, pxlScript, tm)
	if err != nil {
		log.Printf("Error executing PxL script: %s\n", err.Error())
	}
	if err = results.Stream(); err != nil {
		log.Printf("Error executing PxL script: %s\n", err.Error())
	}
}

func (p *pixieMetricsProvider) errors(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Get URL params.
	namespace := ps.ByName("namespace")
	pod := ps.ByName("pod")

	// Pixie refers to pods in the <namespace>/<pod> format.
	pixiePodName := namespace + "/" + pod
	pxlScript := fmt.Sprintf(pxlScript, namespace, pixiePodName, timeWindow)

	// Compute metrics.
	ctx := context.Background()
	p.computeMetrics(ctx, pxlScript)

	// Get metric for pod.
	errorRate := p.podErrorRate[pixiePodName]
	s := fmt.Sprintf("The %s pod(s) has a %2.2f %% error rate.", pod, errorRate)
	log.Println(s)

	// Argo Analysis webhook response needs to requires a JSON response.
	w.Header().Set("Content-Type", "application/json")
	m := map[string]float64{"error_rate": errorRate}
	json.NewEncoder(w).Encode(m)
}

// Implement the TableRecordHandler interface to processes the PxL script output table record-wise.
type podStatsCollector struct {
	podStatsTmp        map[string]float64
	onPodStatsComplete func(stats map[string]float64)
}

func (t *podStatsCollector) HandleInit(ctx context.Context, metadata pxTypes.TableMetadata) error {
	return nil
}

func (t *podStatsCollector) HandleRecord(ctx context.Context, r *pxTypes.Record) error {
	pod := r.GetDatum("pod").String()
	errorRate, ok := r.GetDatum("http_error_rate_in").(*pxTypes.Float64Value)
	if ok {
		t.podStatsTmp[pod] = errorRate.Value()
	}
	return nil
}

func (t *podStatsCollector) HandleDone(ctx context.Context) error {
	t.onPodStatsComplete(t.podStatsTmp)
	return nil
}

// Implement the TableMuxer to route pxl script output tables to the correct handler.
type tableMux struct {
	podStatsCollector  *podStatsCollector
	onPodStatsComplete func(stats map[string]float64)
}

func (s *tableMux) AcceptTable(ctx context.Context, metadata pxTypes.TableMetadata) (pxapi.TableRecordHandler, error) {
	if metadata.Name == "pod_stats" {
		s.podStatsCollector = &podStatsCollector{
			podStatsTmp:        make(map[string]float64),
			onPodStatsComplete: s.onPodStatsComplete,
		}
		return s.podStatsCollector, nil
	}
	return nil, fmt.Errorf("Table %s not found", metadata.Name)
}

func (s *tableMux) GetTable(tableName string) *podStatsCollector {
	if tableName == "pod_stats" {
		return s.podStatsCollector
	}
	return nil
}

func main() {

	log.Println("Starting Pixie metrics server.")

	// Get Pixie API credentials.
	cloudAddr := os.Getenv("PX_CLOUD_ADDR")
	clusterID := os.Getenv("PX_CLUSTER_ID")
	if clusterID == "" {
		log.Fatalln("`PX_CLUSTER_ID` is not set. Did you remember to set the `px-credentials` secret?")
	}
	apiKey := os.Getenv("PX_API_KEY")
	if apiKey == "" {
		log.Fatalln("`PX_API_KEY` is not set. Did you remember to set the `px-credentials` secret?")
	}
	p := newPixieMetricProvider(apiKey, cloudAddr, clusterID)

	router := httprouter.New()
	router.GET("/error-rate/:namespace/:pod", p.errors)
	log.Fatal(http.ListenAndServe(":8080", router))
}
