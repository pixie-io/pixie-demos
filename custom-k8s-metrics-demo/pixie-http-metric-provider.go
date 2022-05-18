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

// The names of the metrics we are creating.
const httpRpsMetricName = "px-http-requests-per-second"
const httpInboundBytesPerSMetricName = "px-http-bytes-recv-per-second"
const httpOutboundBytesPerSMetricName = "px-http-bytes-sent-per-second"
const httpErrorRateMetricName = "px-http-error-rate"
const httpLatencyMsP50 = "px-http-latency-ms-p50"
const httpLatencyMsP90 = "px-http-latency-ms-p90"
const httpLatencyMsP99 = "px-http-latency-ms-p99"

// Map of supported metrics to their corresponding output column name.
var supportedMetrics = map[string]string{
	httpRpsMetricName:               "rps",
	httpErrorRateMetricName:         "error_rate",
	httpInboundBytesPerSMetricName:  "inbound_bytes_per_s",
	httpOutboundBytesPerSMetricName: "outbound_bytes_per_s",
	httpLatencyMsP50:                "latency_ms_p50",
	httpLatencyMsP90:                "latency_ms_p90",
	httpLatencyMsP99:                "latency_ms_p99",
}

// PxL script to compute the metrics. Could be extended to compute additional metrics.
const pxlScript = `import px

# Get list of pods (even non-HTTP)
nanos_per_ms = 1000*1000

df = px.DataFrame(table='process_stats', start_time='-15s')
df.pod = df.ctx['pod']
pods_list = df.groupby('pod').agg()

# Get HTTP events (not all pods will have this)
df = px.DataFrame(table='http_events', start_time='-15s')
df.pod = df.ctx['pod']
df.failure = df.resp_status >= 400
df = df.groupby('pod').agg(
    requests=('latency', px.count),
    error_rate=('failure', px.mean),
    inbound_bytes=('req_body_size', px.sum),
    outbound_bytes=('resp_body_size', px.sum),
    latency_quantiles=('latency', px.quantiles)
)
df.rps = df.requests / 15
df.inbound_bytes_per_s = df.inbound_bytes / 15
df.outbound_bytes_per_s = df.outbound_bytes / 15
df.latency_ms_p50 = px.pluck_float64(df.latency_quantiles, 'p50')/nanos_per_ms
df.latency_ms_p90 = px.pluck_float64(df.latency_quantiles, 'p90')/nanos_per_ms
df.latency_ms_p99 = px.pluck_float64(df.latency_quantiles, 'p99')/nanos_per_ms
df = pods_list.merge(df, how='left', left_on='pod', right_on='pod', suffixes=['', '_x'])
px.display(df[['pod', 'rps', 'error_rate', 'inbound_bytes_per_s', 'outbound_bytes_per_s',
	'latency_ms_p50', 'latency_ms_p90', 'latency_ms_p99']], 'pod_stats')`

// pixieMetricsProvider is a sample implementation of provider.MetricsProvider which computes K8s metrics
// from a PxL script.
type pixieMetricsProvider struct {
	vizierClient         *pxapi.VizierClient
	client               dynamic.Interface
	mapper               apimeta.RESTMapper
	dataMux              sync.Mutex
	podInfo              map[string]map[string]float64
	supportedMetricInfos []provider.CustomMetricInfo
}

func (p *pixieMetricsProvider) computeMetrics(ctx context.Context) {
	tm := &tableMux{
		onPodStatsComplete: func(newStats map[string]map[string]float64) {
			p.dataMux.Lock()
			defer p.dataMux.Unlock()
			p.podInfo = newStats
		},
	}
	log.Println("Refreshing Pixie metrics.")
	results, err := p.vizierClient.ExecuteScript(ctx, pxlScript, tm)
	if err != nil {
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
		<-time.After(15 * time.Second)
	}
}

// NewPixieMetricProvider returns an instance of the Pixie metrics provider.
func NewPixieMetricProvider(vizierClient *pxapi.VizierClient, k8sClient dynamic.Interface, mapper apimeta.RESTMapper) provider.CustomMetricsProvider {
	var supportedMetricInfos []provider.CustomMetricInfo
	for metricName, _ := range supportedMetrics {
		metricInfo := provider.CustomMetricInfo{
			GroupResource: schema.GroupResource{Group: "", Resource: "pods"},
			Metric:        metricName,
			Namespaced:    true,
		}
		supportedMetricInfos = append(supportedMetricInfos, metricInfo)
	}

	provider := &pixieMetricsProvider{
		vizierClient:         vizierClient,
		client:               k8sClient,
		mapper:               mapper,
		podInfo:              make(map[string]map[string]float64),
		supportedMetricInfos: supportedMetricInfos,
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
		Timestamp: metav1.Time{time.Now()},
		Value:     *resource.NewMilliQuantity(int64(value*1000), resource.DecimalSI),
	}, nil
}

// GetMetricByName returns the the pod metric.
func (p *pixieMetricsProvider) GetMetricByName(ctx context.Context, name types.NamespacedName, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValue, error) {
	if _, ok := supportedMetrics[info.Metric]; !ok {
		return nil, provider.NewMetricNotFoundError(info.GroupResource, info.Metric)
	}
	if info.GroupResource.Resource != "pods" {
		return nil, provider.NewMetricNotFoundError(info.GroupResource, info.Metric)
	}

	p.dataMux.Lock()
	defer p.dataMux.Unlock()

	podInfo, ok := p.podInfo[name.String()]
	if !ok {
		return nil, provider.NewMetricNotFoundForError(info.GroupResource, info.Metric, name.String())
	}
	metric, ok := podInfo[info.Metric]
	if !ok {
		return nil, provider.NewMetricNotFoundForError(info.GroupResource, info.Metric, name.String())
	}
	return p.metricFor(metric, name, info)
}

// GetMetricBySelector returns the the pod metric for a pod selector.
func (p *pixieMetricsProvider) GetMetricBySelector(ctx context.Context, namespace string, selector labels.Selector, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValueList, error) {
	names, err := helpers.ListObjectNames(p.mapper, p.client, namespace, selector, info)
	if err != nil {
		return nil, err
	}

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

	return &custom_metrics.MetricValueList{
		Items: res,
	}, nil
}

// ListAllMetrics returns the single metric defined by this provider. Could be extended to return additional Pixie metrics.
func (p *pixieMetricsProvider) ListAllMetrics() []provider.CustomMetricInfo {
	return p.supportedMetricInfos
}

// Implement the TableRecordHandler interface to processes the PxL script output table record-wise.
type podStatsCollector struct {
	podStatsTmp        map[string]map[string]float64
	onPodStatsComplete func(stats map[string]map[string]float64)
}

func (p *podStatsCollector) HandleInit(ctx context.Context, metadata pxTypes.TableMetadata) error {
	p.podStatsTmp = make(map[string]map[string]float64)
	return nil
}

func (p *podStatsCollector) HandleRecord(ctx context.Context, r *pxTypes.Record) error {
	pod := r.GetDatum("pod").String()
	valuesForPod := make(map[string]float64)

	for metricName, metricColumnName := range supportedMetrics {
		metricVal, ok := r.GetDatum(metricColumnName).(*pxTypes.Float64Value)
		if !ok {
			return fmt.Errorf("Metric column %s not found in output table", metricName)
		}
		valuesForPod[metricName] = metricVal.Value()
	}

	p.podStatsTmp[pod] = valuesForPod
	return nil
}

func (p *podStatsCollector) HandleDone(ctx context.Context) error {
	p.onPodStatsComplete(p.podStatsTmp)
	return nil
}

// Implement the TableMuxer to route pxl script output tables to the correct handler.
type tableMux struct {
	podStatsCollector  *podStatsCollector
	onPodStatsComplete func(stats map[string]map[string]float64)
}

func (t *tableMux) AcceptTable(ctx context.Context, metadata pxTypes.TableMetadata) (pxapi.TableRecordHandler, error) {
	if metadata.Name == "pod_stats" {
		t.podStatsCollector = &podStatsCollector{
			podStatsTmp:        make(map[string]map[string]float64),
			onPodStatsComplete: t.onPodStatsComplete,
		}
		return t.podStatsCollector, nil
	}
	return nil, fmt.Errorf("Table %s not found", metadata.Name)
}

func (t *tableMux) GetTable(tableName string) *podStatsCollector {
	if tableName == "pod_stats" {
		return t.podStatsCollector
	}
	return nil
}
