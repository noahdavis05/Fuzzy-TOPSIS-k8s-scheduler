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

type telemetryQuery struct {
	Query      string
	MetricType string // "CPU" or "RAM"
	Stat       string // "Low", "Mean", "High"
}

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

	// make a prometheus client which we can make API calls with
	prometheusClient, err := api.NewClient(api.Config{
		Address: config.PrometheusURL,
	})
	if err != nil {
		fmt.Printf("failed to create Prometheus client: %v\n", err)
	}

	promApi := promv1.NewAPI(prometheusClient)

	refreshedData := requestNodeTelemetry(promApi)

	// the keys are currently the IP adresses of nodes
	// change these back to node.Name
	finalData := make(map[string]types.NodeTelemetryMetrics)
	for _, node := range nodes {
		ip := getNodeAddress(node)
		if ip != "" {
			if metrics, exists := refreshedData[ip]; exists {
				finalData[node.Name] = metrics
			}
		}
	}
	refreshedData = finalData

	//jsonData, _ := json.MarshalIndent(refreshedData, "", "  ")
	//fmt.Println(string(jsonData))
	UpdateCache(refreshedData)
}

func requestNodeTelemetry(promApi promv1.API) map[string]types.NodeTelemetryMetrics {
	queries := []telemetryQuery{
		{`avg by (instance) (100 - (rate(node_cpu_seconds_total{mode="idle"}[5m]) * 100))`, "CPU", "Mean"},
		{`min_over_time((100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[1m])) * 100))[5m:15s])`, "CPU", "Low"},
		{`max_over_time((100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[1m])) * 100))[5m:15s])`, "CPU", "High"},
		{`avg_over_time((100 * (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)))[5m:1m])`, "RAM", "Mean"},
		{`min_over_time((100 * (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)))[5m:1m])`, "RAM", "Low"},
		{`max_over_time((100 * (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)))[5m:1m])`, "RAM", "High"},
	}

	// make a results map
	refreshedData := map[string]types.NodeTelemetryMetrics{}

	// iterate over queries and add them into a result
	for _, query := range queries {
		result := makePromRequest(promApi, query.Query)

		resVector := result.(model.Vector)
		// iterate over the resultVector
		for _, sample := range resVector {
			// we will use the ip address as the key
			instance := string(sample.Metric["instance"])
			ip := strings.Split(instance, ":")[0]

			value := sample.Value

			// check if key exists for this IP, if not make it
			if _, exists := refreshedData[ip]; !exists {
				refreshedData[ip] = types.NodeTelemetryMetrics{}
			}

			metrics := refreshedData[ip]

			// now add the values
			switch query.MetricType {
			case "CPU":
				switch query.Stat {
				case "Mean":
					metrics.CPU.Mean = float64(value)
				case "Low":
					metrics.CPU.Low = float64(value)
				case "High":
					metrics.CPU.High = float64(value)
				}
			case "RAM":
				switch query.Stat {
				case "Mean":
					metrics.RAM.Mean = float64(value)
				case "Low":
					metrics.RAM.Low = float64(value)
				case "High":
					metrics.RAM.High = float64(value)
				}
			}
			refreshedData[ip] = metrics
		}
	}

	//fmt.Printf("Query result: %v\n", resultCPU)
	//fmt.Printf("Query result: %v\n", resultRAM)

	return refreshedData
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
