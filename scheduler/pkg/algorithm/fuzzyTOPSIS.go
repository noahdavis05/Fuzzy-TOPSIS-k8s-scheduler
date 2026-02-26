package algorithm

import (
	"fmt"
	"os"
	"scheduler/pkg/telemetry"
	"text/tabwriter"

	corev1 "k8s.io/api/core/v1"
)

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

	// a list of our columns. E.g. CPU and RAM
	Criteria []string
}

func BuildFuzzyDM(nodes []*corev1.Node) FuzzyDecisionMatrix {
	fuzzyDM := FuzzyDecisionMatrix{
		Data:    make(map[string]map[string]FuzzyNumber),
		Weights: make(map[string]FuzzyNumber),
	}
	fuzzyDM.Criteria = []string{
		"CPU",
		"RAM",
	}

	for _, node := range nodes {
		nodeMetrics, ok := telemetry.GetNodeMetrics(node.Name)
		if !ok {
			fmt.Println("Error getting node metrics")
			panic("Error getting node metrics")
		}

		// make a new row in FuzzyDM for this node
		fuzzyDM.Data[node.Name] = map[string]FuzzyNumber{
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
		}
	}

	return fuzzyDM
}

func DisplayFuzzyDM(fuzzyDM FuzzyDecisionMatrix) {
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
