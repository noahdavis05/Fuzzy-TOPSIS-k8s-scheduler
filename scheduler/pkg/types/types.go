package types

type TelemetryMetric struct {
	Low  float64
	Mean float64
	High float64
}

type NodeTelemetryMetrics struct {
	CPU TelemetryMetric
	RAM TelemetryMetric
}
