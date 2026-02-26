package telemetry

import (
	"context"
	"fmt"
	"scheduler/pkg/config"
	"scheduler/pkg/types"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	v1listers "k8s.io/client-go/listers/core/v1"

	"github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func AutoRefreshTelemetryCache(stopCh <-chan struct{}, interval time.Duration, nodeLister v1listers.NodeLister) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			nodes, err := nodeLister.List(labels.Everything())
			if err != nil {
				fmt.Printf("failed to list nodes: %v\n", err)
				continue
			}
			fmt.Println("Refreshing telemetry cache")
			RefreshTelemetryCache(nodes)
		case <-stopCh:
			fmt.Println("Telemetry refresher finishing")
			return
		}
	}
}

func RefreshTelemetryCache(nodes []*v1.Node) {
	// new data which will replace current cache
	refreshedData := make(map[string]types.NodeTelemetryMetrics)

	// make a prometheus client which we can make API calls with
	prometheusClient, err := api.NewClient(api.Config{
		Address: config.PrometheusURL,
	})
	if err != nil {
		fmt.Printf("failed to create Prometheus client: %v\n", err)
	}

	promApi := promv1.NewAPI(prometheusClient)

	resultCPU, resultRAM := requestNodeTelemetry(promApi)

	cpuVector := resultCPU.(model.Vector)
	ramVector := resultRAM.(model.Vector)

	// map the IPs in the results to actual NodeNames and build
	// refreshed cache map

	// TODO - Can make more efficient
	for _, node := range nodes {
		nodeIP := getNodeAddress(node)

		// now check the cpu and ram and get
		var cpuVal, ramVal float64

		for _, sample := range cpuVector {
			instance := string(sample.Metric["instance"])
			ip := strings.Split(instance, ":")[0]
			if ip == nodeIP {
				cpuVal = float64(sample.Value)
				break
			}
		}

		// Find RAM value
		for _, sample := range ramVector {
			instance := string(sample.Metric["instance"])
			ip := strings.Split(instance, ":")[0]
			if ip == nodeIP {
				ramVal = float64(sample.Value)
				break
			}
		}

		refreshedData[node.Name] = types.NodeTelemetryMetrics{
			CPU: types.TelemetryMetric{
				Low:  0,
				Mean: float32(cpuVal),
				High: 0,
			},
			RAM: types.TelemetryMetric{
				Low:  0,
				Mean: float32(ramVal),
				High: 0,
			},
		}
	}
	fmt.Printf("Refreshed telemetry data: %+v\n", refreshedData)
	UpdateCache(refreshedData)
}

func requestNodeTelemetry(promApi promv1.API) (model.Value, model.Value) {
	query := `sum by (instance) (rate(node_cpu_seconds_total{mode!="idle"}[5m]))`

	resultCPU := makePromRequest(promApi, query)

	query = `sum by (instance) (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes)`

	resultRAM := makePromRequest(promApi, query)

	//fmt.Printf("Query result: %v\n", resultCPU)
	//fmt.Printf("Query result: %v\n", resultRAM)

	return resultCPU, resultRAM
}

func makePromRequest(promAPI promv1.API, query string) model.Value {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, warnings, err := promAPI.Query(ctx, query, time.Now())
	if err != nil {
		fmt.Printf("Prometheus query error: %v", err)
	}
	if len(warnings) > 0 {
		fmt.Println("Prometheus query warnings:", warnings)
	}

	return result
}

func getNodeAddress(node *v1.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == v1.NodeInternalIP {
			return addr.Address
		}
	}
	return ""
}
