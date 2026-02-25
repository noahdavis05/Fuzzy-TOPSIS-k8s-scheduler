package telemetry

import (
	"scheduler/pkg/types"

	v1 "k8s.io/api/core/v1"
)

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
