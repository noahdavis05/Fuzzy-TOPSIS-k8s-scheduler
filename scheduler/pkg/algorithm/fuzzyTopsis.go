package algorithm

import (
	"encoding/json"
	"fmt"
	"math"
	"scheduler/pkg/types"
)

func SelectNode(fuzzyDM types.FuzzyDecisionMatrix) string {
	// all our values in fuzzyDM are percentages e.g. between 0 and 100
	// therefore already normalised/on same scale
	weightNodes(&fuzzyDM)
	weightIdeals(&fuzzyDM)
	DisplayFuzzyDM(fuzzyDM)

	nodeScores := scoreNodes(fuzzyDM)

	fmt.Println("Node scores:")
	data, _ := json.MarshalIndent(nodeScores, "", "  ")
	fmt.Println(string(data))

	// now get the key of the node with the highest value
	nodeName := ""
	maxScore := -math.Inf(1)
	for node, score := range nodeScores {
		if score > maxScore {
			maxScore = score
			nodeName = node
		}
	}
	//
	return nodeName
}

func scoreNodes(fuzzyDM types.FuzzyDecisionMatrix) map[string]float64 {
	// iterate over the fuzzyDM and score each node
	nodeScores := map[string]float64{}

	// iterate over the nodes and score their positive and negative distances
	for node, criterion := range fuzzyDM.Data {
		// iterate over all criteria in each node
		negativeDists := float64(0)
		positiveDists := float64(0)
		for criteria, value := range criterion {
			fuzzyNum := value
			positiveIdeal := fuzzyDM.PositiveIdeals[criteria]
			negativeIdeal := fuzzyDM.NegativeIdeals[criteria]
			positiveDists += calculateDistance(fuzzyNum, positiveIdeal)
			negativeDists += calculateDistance(fuzzyNum, negativeIdeal)
		}
		// now with positive and negative dists
		nodeScore := negativeDists / (negativeDists + positiveDists)
		nodeScores[node] = nodeScore
	}
	return nodeScores
}

func calculateDistance(fuzzyNum types.FuzzyNumber, fuzzyIdeal types.FuzzyNumber) float64 {
	dist1 := (fuzzyNum.A - fuzzyIdeal.A) * (fuzzyNum.A - fuzzyIdeal.A)
	dist2 := (fuzzyNum.B - fuzzyIdeal.B) * (fuzzyNum.B - fuzzyIdeal.B)
	dist3 := (fuzzyNum.C - fuzzyIdeal.C) * (fuzzyNum.C - fuzzyIdeal.C)

	totalSquaredDistances := (dist1 + dist2 + dist3) / 3

	return math.Sqrt(totalSquaredDistances)
}

func weightNodes(fuzzyDM *types.FuzzyDecisionMatrix) {
	// the desicion matrix is passed as pointer so doesn't need to be changed
	for k, v := range fuzzyDM.Data {
		for key, value := range v {
			// key is field e.g. CPU
			// value is the FuzzyNumber we need to update
			weights := fuzzyDM.Weights[key]
			weightedFuzzyNum := types.FuzzyNumber{
				A: value.A * weights.A,
				B: value.B * weights.B,
				C: value.C * weights.C,
			}
			fuzzyDM.Data[k][key] = weightedFuzzyNum
		}
	}
}

func weightIdeals(fuzzyDM *types.FuzzyDecisionMatrix) {
	for key, value := range fuzzyDM.PositiveIdeals {
		weights := fuzzyDM.Weights[key]
		weightedFuzzyNum := types.FuzzyNumber{
			A: value.A * weights.A,
			B: value.B * weights.B,
			C: value.C * weights.C,
		}
		fuzzyDM.PositiveIdeals[key] = weightedFuzzyNum
	}

	for key, value := range fuzzyDM.NegativeIdeals {
		weights := fuzzyDM.Weights[key]
		weightedFuzzyNum := types.FuzzyNumber{
			A: value.A * weights.A,
			B: value.B * weights.B,
			C: value.C * weights.C,
		}
		fuzzyDM.NegativeIdeals[key] = weightedFuzzyNum
	}
}
