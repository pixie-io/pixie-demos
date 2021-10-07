/*
Copyright 2018 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"log"
	"os"

	"k8s.io/apimachinery/pkg/util/wait"
	basecmd "sigs.k8s.io/custom-metrics-apiserver/pkg/cmd"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"

	"px.dev/pxapi"
)

// Adapted from the example in this repo: https://github.com/kubernetes-sigs/custom-metrics-apiserver

type pixieAdapter struct {
	basecmd.AdapterBase
	Message string
}

func (a *pixieAdapter) makeProviderOrDie(clusterID string, apiKey string, cloudAddr string) provider.CustomMetricsProvider {
	ctx := context.Background()
	pixieClient, err := pxapi.NewClient(ctx, pxapi.WithAPIKey(apiKey), pxapi.WithCloudAddr(cloudAddr))
	if err != nil {
		log.Fatalln(err.Error())
	}
	vz, err := pixieClient.NewVizierClient(ctx, clusterID)
	if err != nil {
		log.Fatalln(err.Error())
	}

	client, err := a.DynamicClient()
	if err != nil {
		log.Fatalf("unable to construct dynamic client: %v", err)
	}

	mapper, err := a.RESTMapper()
	if err != nil {
		log.Fatalf("unable to construct discovery REST mapper: %v", err)
	}

	return NewPixieMetricProvider(vz, client, mapper)
}

func main() {
	cmd := &pixieAdapter{
		Message: "Starting Pixie custom metrics adapter",
	}
	cmd.Flags().Parse(os.Args)

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

	testProvider := cmd.makeProviderOrDie(clusterID, apiKey, cloudAddr)
	cmd.WithCustomMetrics(testProvider)

	log.Println(cmd.Message)
	if err := cmd.Run(wait.NeverStop); err != nil {
		log.Fatalf("unable to run custom metrics adapter: %v", err)
	}
}
