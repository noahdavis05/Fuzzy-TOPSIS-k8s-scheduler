package cluster

import (
	"fmt"
	"scheduler/pkg/types"

	v1 "k8s.io/api/core/v1"
)

// funtion builds a clusterInfo struct on each nodes' limits
// we use Allocatable which is all the total resources available
// for pods to be scheduled with. E.g. full total resources
// minus system overhead
func CreateClusterInfo(nodes []*v1.Node) types.ClusterInfo {
	clusterInfo := types.ClusterInfo{
		CPULimits: make(map[string]int64),
		RAMLimits: make(map[string]int64),
	}

	for _, node := range nodes {
		cpuAllocatable := node.Status.Allocatable.Cpu().MilliValue()
		ramAllocatable := node.Status.Allocatable.Memory().Value()
		clusterInfo.CPULimits[node.Name] = cpuAllocatable
		clusterInfo.RAMLimits[node.Name] = ramAllocatable
	}

	fmt.Println(clusterInfo)

	return clusterInfo
}
