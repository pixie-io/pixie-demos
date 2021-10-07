package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	apierr "k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/metrics/pkg/apis/custom_metrics"

	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider/helpers"

	"px.dev/pxapi"
	pxTypes "px.dev/pxapi/types"
)

// Adapted from the example in this repo: https://github.com/kubernetes-sigs/custom-metrics-apiserver

// The name of the metric we are creating.
const httpRpsMetricName = "px-http-requests-per-second"

// PxL script to compute the metrics. Could be extended to compute additional metrics.
const pxlScript = `import px

# Get list of pods (even non-HTTP)
df = px.DataFrame(table='process_stats', start_time='-30s')
df.pod = df.ctx['pod']
pods_list = df.groupby('pod').agg()

# Get HTTP events (not all pods will have this)
df = px.DataFrame(table='http_events', start_time='-30s')
df.pod = df.ctx['pod']
df = df.groupby('pod').agg(
    requests=('latency', px.count)
)
df.rps = df.requests / 30;
df = pods_list.merge(df, how='left', left_on='pod', right_on='pod', suffixes=['', '_x'])
px.display(df[['pod', 'rps']], 'pod_stats')`

// pixieMetricsProvider is a sample implementation of provider.MetricsProvider which computes K8s metrics
// from a PxL script.
type pixieMetricsProvider struct {
	vizierClient *pxapi.VizierClient
	client dynamic.Interface
	mapper apimeta.RESTMapper
	dataMux sync.Mutex
	podRequestsPerS map[string]float64
}

func (p *pixieMetricsProvider) computeMetrics(ctx context.Context) {
	tm := &tableMux{
		onPodStatsComplete: func(newStats map[string]float64) {
			p.dataMux.Lock()
			defer p.dataMux.Unlock()
			p.podRequestsPerS = newStats
		},
	}
	log.Println("Refreshing Pixie metrics.")
	results, err := p.vizierClient.ExecuteScript(ctx, pxlScript, tm)
	if err != nil  {
		log.Printf("Error executing PxL script: %s\n", err.Error())
	}
	if err = results.Stream(); err != nil {
		log.Printf("Error executing PxL script: %s\n", err.Error())
	}
}

func (p *pixieMetricsProvider) runMetricsLoop() {
	ctx := context.Background()
  for {
		p.computeMetrics(ctx)
    <-time.After(30 * time.Second)
  }
}

// NewPixieMetricProvider returns an instance of the Pixie metrics provider.
func NewPixieMetricProvider(vizierClient *pxapi.VizierClient, k8sClient dynamic.Interface, mapper apimeta.RESTMapper) provider.CustomMetricsProvider {
	provider := &pixieMetricsProvider{
		vizierClient:    vizierClient,
		client:          k8sClient,
		mapper:          mapper,
		podRequestsPerS:  make(map[string]float64),
	}
	go provider.runMetricsLoop()
	return provider
}

func (p *pixieMetricsProvider) metricFor(value float64, name types.NamespacedName, info provider.CustomMetricInfo) (*custom_metrics.MetricValue, error) {
	// construct a reference referring to the described object
	objRef, err := helpers.ReferenceFor(p.mapper, name, info)
	if err != nil {
			return nil, err
	}

	return &custom_metrics.MetricValue{
			DescribedObject: objRef,
			Metric: custom_metrics.MetricIdentifier{
				Name: info.Metric,
			},
			Timestamp:       metav1.Time{time.Now()},
			Value:           *resource.NewMilliQuantity(int64(value*1000), resource.DecimalSI),
	}, nil
}

// GetMetricByName returns the the pod RPS metric.
func (p *pixieMetricsProvider) GetMetricByName(ctx context.Context, name types.NamespacedName, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValue, error) {
	if (info.Metric != httpRpsMetricName || info.GroupResource.Resource != "pods") {
		return nil, provider.NewMetricNotFoundError(info.GroupResource, info.Metric)
	}
	rps, ok := p.podRequestsPerS[name.String()]
	if !ok {
		return nil, provider.NewMetricNotFoundForError(info.GroupResource, info.Metric, name.String())
	}
	return p.metricFor(rps, name, info)
}

// GetMetricBySelector returns the the pod RPS metric for a pod selector.
func (p *pixieMetricsProvider) GetMetricBySelector(ctx context.Context, namespace string, selector labels.Selector, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValueList, error) {
	names, err := helpers.ListObjectNames(p.mapper, p.client, namespace, selector, info)
	if err != nil {
		return nil, err
	}

	fmt.Printf("GetMetricBySelector names: %+s", names)

	res := make([]custom_metrics.MetricValue, 0, len(names))
	for _, name := range names {
		namespacedName := types.NamespacedName{Name: name, Namespace: namespace}
		value, err := p.GetMetricByName(ctx, namespacedName, info, metricSelector)
		if err != nil || value == nil {
			if apierr.IsNotFound(err) {
				continue
			}
			return nil, err
		}
		res = append(res, *value)
	}

	fmt.Printf("GetMetricBySelector Value list: %+s", res)

	return &custom_metrics.MetricValueList{
		Items: res,
	}, nil
}

// ListAllMetrics returns the single metric defined by this provider. Could be extended to return additional Pixie metrics.
func (p *pixieMetricsProvider) ListAllMetrics() []provider.CustomMetricInfo {
	return []provider.CustomMetricInfo{
			{
					GroupResource: schema.GroupResource{Group: "", Resource: "pods"},
					Metric:        httpRpsMetricName,
					Namespaced:    true,
			},
	}
}

// Implement the TableRecordHandler interface to processes the PxL script output table record-wise.
type podStatsCollector struct {
	podStatsTmp map[string]float64
	onPodStatsComplete func(stats map[string]float64)
}

func (t *podStatsCollector) HandleInit(ctx context.Context, metadata pxTypes.TableMetadata) error {
	return nil
}

func (t *podStatsCollector) HandleRecord(ctx context.Context, r *pxTypes.Record) error {
	pod := r.GetDatum("pod").String()
	rps, ok := r.GetDatum("rps").(*pxTypes.Float64Value)
	if ok {
		t.podStatsTmp[pod] = rps.Value()
	}
	return nil
}

func (t *podStatsCollector) HandleDone(ctx context.Context) error {
	t.onPodStatsComplete(t.podStatsTmp)
	return nil
}

// Implement the TableMuxer to route pxl script output tables to the correct handler.
type tableMux struct {
	podStatsCollector *podStatsCollector
	onPodStatsComplete func(stats map[string]float64)
}

func (s *tableMux) AcceptTable(ctx context.Context, metadata pxTypes.TableMetadata) (pxapi.TableRecordHandler, error) {
	if metadata.Name == "pod_stats" {
		s.podStatsCollector = &podStatsCollector{
			podStatsTmp: make(map[string]float64),
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
