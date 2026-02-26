package scheduler

import (
	"context"
	"fmt"
	"scheduler/pkg/algorithm"
	"scheduler/pkg/telemetry"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	v1listers "k8s.io/client-go/listers/core/v1"
)

func SchedulePod(client *kubernetes.Clientset, pod *corev1.Pod, nodeLister v1listers.NodeLister) {
	nodes, err := nodeLister.List(labels.Everything())
	if err != nil {
		fmt.Printf("failed to list nodes: %v\n", err)
	}

	fmt.Println("Loading telemetry data")
	telemetry.RefreshTelemetryCache(nodes)

	// create fuzzy decision matrix and print in terminal for debugging
	fuzzyDM := algorithm.BuildFuzzyDM(nodes)
	algorithm.DisplayFuzzyDM(fuzzyDM)

	selectedNode := nodes[2]

	// get the telemetry data from this node
	nodeTelemetry, ok := telemetry.GetNodeMetrics(selectedNode.Name)
	if !ok {
		panic("Error getting telemetry")
	}
	fmt.Printf("Node mean CPU %f, Node mean RAM %f\n", nodeTelemetry.CPU.Mean, nodeTelemetry.RAM.High)
	bindPod(client, pod, selectedNode)
}

// Bind a pod to a Node
func bindPod(client *kubernetes.Clientset, pod *corev1.Pod, node *corev1.Node) {
	binding := &v1.Binding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			UID:       pod.UID,
		},
		Target: v1.ObjectReference{
			Kind: "Node",
			Name: node.Name,
		},
	}

	err := client.CoreV1().Pods(pod.Namespace).Bind(context.TODO(), binding, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("Error updating the nodename")
	}

}
