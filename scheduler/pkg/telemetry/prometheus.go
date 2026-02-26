package telemetry

import (
	"fmt"
	"scheduler/pkg/types"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	v1listers "k8s.io/client-go/listers/core/v1"
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
	refreshedData := make(map[string]types.NodeTelemetryMetrics)

	for _, node := range nodes {
		_ = node
		// make individual request
		result := requestNodeTelemetry(node.Name)

		// add result to refreshedData
		refreshedData[node.Name] = result
	}

	UpdateCache(refreshedData)
}

func requestNodeTelemetry(nodeName string) types.NodeTelemetryMetrics {
	return types.NodeTelemetryMetrics{
		CPU: types.TelemetryMetric{
			Low:  0,
			Mean: 50,
			High: 10,
		},
		RAM: types.TelemetryMetric{
			Low:  0,
			Mean: 50,
			High: 10,
		},
	}
}
