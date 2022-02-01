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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/julienschmidt/httprouter"
	"px.dev/pxapi"
	pxTypes "px.dev/pxapi/types"
)

func (p *pixieMetricsProvider) podRPS(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	// Read PxL script from file.
	b, err := ioutil.ReadFile("svc_error_rate.pxl")
	if err != nil {
		panic(err)
	}
	pxlScript := string(b)

	// Compute metrics
	ctx := context.Background()
	p.computeMetrics(ctx, pxlScript)
	//log.Println("New Pixie metrics: ", p.podRequestsPerS)

	// Get URL params.
	namespace := ps.ByName("namespace")
	service := ps.ByName("service")
	// Pixie refers to pods and services in the <namespace>/<pod,service> format.
	pixieServiceName := namespace + "/" + service

	// Get metric for pod.
	errorRate := p.podRequestsPerS[pixieServiceName]

	log.Println(service, " has ", errorRate, " error rate.")

	// Argo Analysis webhook response must return JSON content.
	w.Header().Set("Content-Type", "application/json")
	m := map[string]float64{"error_rate": errorRate}
	json.NewEncoder(w).Encode(m)
}

type pixieMetricsProvider struct {
	vizierClient    *pxapi.VizierClient
	dataMux         sync.Mutex
	podRequestsPerS map[string]float64
}

func (p *pixieMetricsProvider) computeMetrics(ctx context.Context, pxlScript string) {
	tm := &tableMux{
		onPodStatsComplete: func(newStats map[string]float64) {
			p.dataMux.Lock()
			defer p.dataMux.Unlock()
			p.podRequestsPerS = newStats
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

// NewPixieMetricProvider returns an instance of the Pixie metrics provider.
func NewPixieMetricProvider(apiKey string, cloudAddr string, clusterID string) *pixieMetricsProvider {

	// Create Pixie client.
	log.Println("Creating Pixie client.")
	ctx := context.Background()
	pixieClient, err := pxapi.NewClient(ctx, pxapi.WithAPIKey(apiKey), pxapi.WithCloudAddr(cloudAddr))
	if err != nil {
		log.Fatalln(err.Error())
	}
	vz, err := pixieClient.NewVizierClient(ctx, clusterID)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Create pixieMetricsProvider.
	provider := &pixieMetricsProvider{
		vizierClient:    vz,
		podRequestsPerS: make(map[string]float64),
	}

	return provider
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
	service := r.GetDatum("service").String()
	errorRate, ok := r.GetDatum("http_error_rate_in").(*pxTypes.Float64Value)
	if ok {
		t.podStatsTmp[service] = errorRate.Value()
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
	if cloudAddr == "" {
		log.Fatalln("`PX_CLOUD_ADDR` is not set.")
	}
	clusterID := os.Getenv("PX_CLUSTER_ID")
	if clusterID == "" {
		log.Fatalln("`PX_CLUSTER_ID` is not set. Did you remember to set the `px-credentials` secret?")
	}
	apiKey := os.Getenv("PX_API_KEY")
	if apiKey == "" {
		log.Fatalln("`PX_API_KEY` is not set. Did you remember to set the `px-credentials` secret?")
	}

	p := NewPixieMetricProvider(apiKey, cloudAddr, clusterID)

	router := httprouter.New()
	router.GET("/podrps/:namespace/:service", p.podRPS)
	log.Fatal(http.ListenAndServe(":8080", router))
}
