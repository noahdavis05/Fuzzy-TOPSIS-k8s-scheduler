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

type FuzzyNumber struct {
	A float64
	B float64
	C float64
}

type FuzzyDecisionMatrix struct {
	// e.g. Data["node1"]["CPU"] = fuzzy number
	Data map[string]map[string]FuzzyNumber

	// e.g. Weights["CPU"] = fuzzy number
	Weights map[string]FuzzyNumber

	// positive ideal solutions - e.g. what we want our nodes to look like
	// e.g. PositiveIdeals["CPU"] = (75,75,75)
	PositiveIdeals map[string]FuzzyNumber

	// negative ideal solutions - e.g. what we don't want our nodes to look like
	// e.g. NegativeIdeals["CPU"] = (100,100,100)
	NegativeIdeals map[string]FuzzyNumber

	// a list of our columns. E.g. CPU and RAM
	Criteria []string
}
