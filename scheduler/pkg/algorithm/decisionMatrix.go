package algorithm

import (
	"fmt"
	"os"
	"scheduler/pkg/config"
	"scheduler/pkg/telemetry"
	"scheduler/pkg/types"
	"text/tabwriter"

	corev1 "k8s.io/api/core/v1"
)

func BuildFuzzyDM(nodes []*corev1.Node) types.FuzzyDecisionMatrix {
	fuzzyDM := types.FuzzyDecisionMatrix{
		Data: make(map[string]map[string]types.FuzzyNumber),
	}
	fuzzyDM.Criteria = []string{
		"CPU",
		"RAM",
		"CPU RANGE",
		"RAM RANGE",
	}

	// These are the weights used as part of TOPSIS
	fuzzyDM.Weights = map[string]types.FuzzyNumber{
		"CPU":       config.CPUWeights,
		"RAM":       config.RAMWeights,
		"CPU RANGE": config.CPURangeWeights,
		"RAM RANGE": config.RAMRangeWeights,
	}

	// set the Ideal Positives and Ideal Negatives
	fuzzyDM.PositiveIdeals = map[string]types.FuzzyNumber{
		"CPU":       config.PosCPUIdeal,
		"RAM":       config.PosRAMIdeal,
		"CPU RANGE": config.PosCPURangeIdeal,
		"RAM RANGE": config.PosRAMRangeIdeal,
	}

	fuzzyDM.NegativeIdeals = map[string]types.FuzzyNumber{
		"CPU":       config.NegCPUIdeal,
		"RAM":       config.NegRAMIdeal,
		"CPU RANGE": config.NegCPURangeIdeal,
		"RAM RANGE": config.NegRAMRangeIdeal,
	}

	for _, node := range nodes {
		nodeMetrics, ok := telemetry.GetNodeMetrics(node.Name)
		if !ok {
			fmt.Println("Error getting node metrics")
			panic("Error getting node metrics")
		}

		// make a new row in FuzzyDM for this node
		fuzzyDM.Data[node.Name] = map[string]types.FuzzyNumber{
			"CPU": {
				A: nodeMetrics.CPU.Low,
				B: nodeMetrics.CPU.Mean,
				C: nodeMetrics.CPU.High,
			},
			"RAM": {
				A: nodeMetrics.RAM.Low,
				B: nodeMetrics.RAM.Mean,
				C: nodeMetrics.RAM.High,
			},
			"CPU RANGE": {
				A: 0,
				B: nodeMetrics.CPU.High - nodeMetrics.CPU.Low,
				C: nodeMetrics.CPU.High - nodeMetrics.CPU.Low,
			},
			"RAM RANGE": {
				A: 0,
				B: nodeMetrics.RAM.High - nodeMetrics.RAM.Low,
				C: nodeMetrics.RAM.High - nodeMetrics.RAM.Low,
			},
		}
	}

	return fuzzyDM
}

func DisplayFuzzyDM(fuzzyDM types.FuzzyDecisionMatrix) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

	// build title line
	fmt.Fprint(w, "Node\t")
	for _, criterion := range fuzzyDM.Criteria {
		fmt.Fprintf(w, " %s (a, b, c)\t", criterion)
	}
	fmt.Fprintln(w)

	// add seperator
	fmt.Fprint(w, "---\t")
	for range fuzzyDM.Criteria {
		fmt.Fprint(w, "-----------\t")
	}
	fmt.Fprintln(w)

	// print rows
	for nodeName, metrics := range fuzzyDM.Data {
		fmt.Fprintf(w, "%s\t", nodeName)
		for _, criterion := range fuzzyDM.Criteria {
			f := metrics[criterion]
			fmt.Fprintf(w, " (%.2f, %.2f, %.2f)\t", f.A, f.B, f.C)
		}
		fmt.Fprintln(w)
	}

	w.Flush()
	fmt.Println()
}
