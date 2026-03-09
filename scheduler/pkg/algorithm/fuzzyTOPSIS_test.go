package algorithm

import (
	"fmt"
	"scheduler/pkg/dashboard"
	"scheduler/pkg/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

// These tests will mock the whole decision matrix for numerous scenarios
// and ensure we get the node selected that we want
func TestTOPSISRankings(t *testing.T) {
	// set constant vars for all these tests
	// cluster limits - always three nodes, all have the same limits
	clusterLimits := types.ClusterInfo{
		CPULimits: map[string]int64{
			"Node1": 2000,
			"Node2": 2000,
			"Node3": 2000,
		},
		RAMLimits: map[string]int64{
			"Node1": 3809775616,
			"Node2": 3809775616,
			"Node3": 3809775616,
		},
	}

	// podRequests are just needed for the filtering
	// not too important in this test, just keep them small
	// so nodes don't need to be filtered unless usage already
	// very high. Filtering is tested thoroughly in another
	// test
	podRequest := types.PodRequest{
		CPU: 100,
		RAM: 250,
	}
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
		{
			// Test node 1 gets filtered out
			name: "Test filtering",
			decisionMatrix: map[string]map[string]types.FuzzyNumber{
				"Node1": {
					"CPU": {
						A: 95,
						B: 97,
						C: 100,
					},
					"RAM": {
						A: 95,
						B: 97,
						C: 100,
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
				"Node2": {
					"CPU": {
						A: 55,
						B: 60,
						C: 65,
					},
					"RAM": {
						A: 55,
						B: 60,
						C: 65,
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
						A: 10,
						B: 20,
						C: 30,
					},
					"RAM": {
						A: 10,
						B: 20,
						C: 30,
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
			expectedNodeName: "Node2",
		},
		{
			// Test when nodes are all heavily packed
			// and equal ranges.
			name: "Test dangerzone",
			decisionMatrix: map[string]map[string]types.FuzzyNumber{
				"Node1": {
					"CPU": {
						A: 77,
						B: 82,
						C: 87,
					},
					"RAM": {
						A: 77,
						B: 82,
						C: 87,
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
						A: 72,
						B: 77,
						C: 82,
					},
					"RAM": {
						A: 72,
						B: 77,
						C: 82,
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
						A: 67,
						B: 72,
						C: 77,
					},
					"RAM": {
						A: 67,
						B: 72,
						C: 77,
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
			expectedNodeName: "Node3",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n\nRunning Test Case ==> %v\n", tc.name)
			// create fuzzyDM
			fuzzyDM := buildTestingDM()
			fuzzyDM.Data = tc.decisionMatrix
			FilterNodes(&fuzzyDM, podRequest, clusterLimits)

			// now run the node selection#
			dash := dashboard.PodScheduledMessage{}
			outputNodeName := selectNode(fuzzyDM, &dash, false)
			assert.Equal(t, outputNodeName, tc.expectedNodeName)

			if outputNodeName == tc.expectedNodeName {
				fmt.Printf("\033[32mTest passed: %v\033[0m\n", tc.name)
			} else {
				fmt.Printf("\033[31mTest failed: %v\033[0m\n", tc.name)
			}
		})
	}
}

// The above tests do filter but just assume
// filtering works. These tests check filtering
// in more detail.
func TestFiltering(t *testing.T) {
	// cluster Limits are the same for all tests
	// we are just testing one node per test
	// therefore cluster limits only has one node
	// this node will just be called Node1
	nodeName := "Node1"

	clusterNodes := types.ClusterInfo{
		CPULimits: map[string]int64{
			"Node1": 2000, // 2 cores
		},
		RAMLimits: map[string]int64{
			"Node1": 3809775616, // 3.8 GB
		},
	}

	// make scenarios for Node1's state e.g. under heavy load and under low load
	// 75% load on all 10% range
	scneario1 := map[string]map[string]types.FuzzyNumber{
		nodeName: {
			"CPU": types.FuzzyNumber{
				A: 70,
				B: 75,
				C: 80,
			},
			"RAM": types.FuzzyNumber{
				A: 70,
				B: 75,
				C: 80,
			},
		},
	}
	tests := []struct {
		name           string
		podRequest     types.PodRequest
		nodeTelemetry  map[string]map[string]types.FuzzyNumber
		expectedResult bool
	}{
		{
			name: "Test node doesn't get filtered",
			podRequest: types.PodRequest{
				CPU: 200,
				RAM: 500000,
			},
			nodeTelemetry:  scneario1,
			expectedResult: false,
		},
		{
			name: "Test pod requests too much CPU",
			podRequest: types.PodRequest{
				CPU: 250,
				RAM: 500000,
			},
			nodeTelemetry:  scneario1,
			expectedResult: true,
		},
		{
			name: "Test pod requests too much RAM",
			podRequest: types.PodRequest{
				CPU: 200,
				RAM: 381000000,
			},
			nodeTelemetry:  scneario1,
			expectedResult: true,
		},
		{
			name: "Test pod requests too much RAM and CPU",
			podRequest: types.PodRequest{
				CPU: 250,
				RAM: 381000000,
			},
			nodeTelemetry:  scneario1,
			expectedResult: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// create fuzzyDM
			fmt.Printf("\n\nRunning Test Case ==> %v\n", tc.name)
			fuzzyDM := buildTestingDM()
			fuzzyDM.Data = tc.nodeTelemetry
			res := filterNode(&fuzzyDM, nodeName, tc.podRequest, clusterNodes)
			// proper assertion
			assert.Equal(t, res, tc.expectedResult)

			// another check to print out a success message in green or error in red
			if res == tc.expectedResult {
				fmt.Printf("\033[32mTest passed: %v\033[0m\n", tc.name)
			} else {
				fmt.Printf("\033[31mTest failed: %v\033[0m\n", tc.name)
			}
		})
	}
}
