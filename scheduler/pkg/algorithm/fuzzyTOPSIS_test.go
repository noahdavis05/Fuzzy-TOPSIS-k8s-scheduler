package algorithm

import (
	"fmt"
	"scheduler/pkg/types"
	"testing"
)

// These tests will mock the whole decision matrix for numerous scenarios
// and ensure we get the node selected that we want
func TestTOPSISRankings(t *testing.T) {
	tests := []struct {
		name             string
		decisionMatrix   map[string]map[string]types.FuzzyNumber
		expectedNodeName string
	}{
		{
			// In this scenario node 1 is perfect - e.g. matches ideals
			// node 2 is still quite good, but bigger range
			// node 3 is quite unstable
			name: "Node 1 'perfect' choice",
			decisionMatrix: map[string]map[string]types.FuzzyNumber{
				"Node1": {
					"CPU": {
						A: 74,
						B: 75,
						C: 76,
					},
					"RAM": {
						A: 74,
						B: 75,
						C: 76,
					},
					"CPU RANGE": {
						A: 0,
						B: 1,
						C: 1,
					},
					"RAM RANGE": {
						A: 0,
						B: 1,
						C: 1,
					},
				},
				"Node2": {
					"CPU": {
						A: 60,
						B: 65,
						C: 70,
					},
					"RAM": {
						A: 60,
						B: 65,
						C: 70,
					},
					"CPU RANGE": {
						A: 0,
						B: 10,
						C: 10,
					},
					"RAM RANGE": {
						A: 0,
						B: 10,
						C: 10,
					},
				},
				"Node3": {
					"CPU": {
						A: 60,
						B: 70,
						C: 80,
					},
					"RAM": {
						A: 60,
						B: 70,
						C: 80,
					},
					"CPU RANGE": {
						A: 0,
						B: 20,
						C: 20,
					},
					"RAM RANGE": {
						A: 0,
						B: 20,
						C: 20,
					},
				},
			},
			expectedNodeName: "Node1",
		},
		{
			// In this scenario node 1 is packed and nodes 2 and 3 are empty
			// all have similar stability
			// should expect node 1 to be chosen
			name: "Test Bin Pack",
			decisionMatrix: map[string]map[string]types.FuzzyNumber{
				"Node1": {
					"CPU": {
						A: 70,
						B: 75,
						C: 80,
					},
					"RAM": {
						A: 65,
						B: 70,
						C: 75,
					},
					"CPU RANGE": {
						A: 0,
						B: 10,
						C: 10,
					},
					"RAM RANGE": {
						A: 0,
						B: 10,
						C: 10,
					},
				},
				"Node2": {
					"CPU": {
						A: 10,
						B: 15,
						C: 20,
					},
					"RAM": {
						A: 20,
						B: 25,
						C: 30,
					},
					"CPU RANGE": {
						A: 0,
						B: 10,
						C: 10,
					},
					"RAM RANGE": {
						A: 0,
						B: 10,
						C: 10,
					},
				},
				"Node3": {
					"CPU": {
						A: 15,
						B: 20,
						C: 25,
					},
					"RAM": {
						A: 15,
						B: 20,
						C: 25,
					},
					"CPU RANGE": {
						A: 0,
						B: 10,
						C: 10,
					},
					"RAM RANGE": {
						A: 0,
						B: 10,
						C: 10,
					},
				},
			},
			expectedNodeName: "Node1",
		},
		{
			// In this scenario node 1 is packed and very unstable
			// nodes 2 and 3 are quite empty but stable
			// should expect node 2 to be chosen
			// node 2 will have the same stability as node 3 but
			// slightly higher utilisation
			name: "Test Bin Pack",
			decisionMatrix: map[string]map[string]types.FuzzyNumber{
				"Node1": {
					"CPU": {
						A: 70,
						B: 80,
						C: 90,
					},
					"RAM": {
						A: 65,
						B: 75,
						C: 85,
					},
					"CPU RANGE": {
						A: 0,
						B: 20,
						C: 20,
					},
					"RAM RANGE": {
						A: 0,
						B: 20,
						C: 20,
					},
				},
				"Node2": {
					"CPU": {
						A: 15,
						B: 17,
						C: 20,
					},
					"RAM": {
						A: 20,
						B: 22,
						C: 25,
					},
					"CPU RANGE": {
						A: 0,
						B: 5,
						C: 5,
					},
					"RAM RANGE": {
						A: 0,
						B: 5,
						C: 5,
					},
				},
				"Node3": {
					"CPU": {
						A: 10,
						B: 12,
						C: 15,
					},
					"RAM": {
						A: 15,
						B: 17,
						C: 20,
					},
					"CPU RANGE": {
						A: 0,
						B: 5,
						C: 5,
					},
					"RAM RANGE": {
						A: 0,
						B: 5,
						C: 5,
					},
				},
			},
			expectedNodeName: "Node2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// create fuzzyDM
			fuzzyDM := buildTestingDM()
			fuzzyDM.Data = tc.decisionMatrix

			// now run the node selection
			outputNodeName := SelectNode(fuzzyDM)
			fmt.Println(outputNodeName)
		})
	}

}
