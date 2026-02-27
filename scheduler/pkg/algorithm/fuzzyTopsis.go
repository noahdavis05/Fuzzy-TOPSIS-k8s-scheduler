package algorithm

import "scheduler/pkg/types"

func SelectNode(fuzzyDM types.FuzzyDecisionMatrix) (nodeName string) {
	// all our values in fuzzyDM are percentages e.g. between 0 and 100
	// therefore already normalised/on same scale
	weightNodes(&fuzzyDM)
	weightIdeals(&fuzzyDM)
	DisplayFuzzyDM(fuzzyDM)
	//
	return "worker"
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
