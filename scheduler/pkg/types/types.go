package types

type TelemetryMetric struct {
	Low  float32
	Mean float32
	High float32
}

type NodeTelemetryMetrics struct {
	CPU TelemetryMetric
	RAM TelemetryMetric
}
