package main

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"scheduler/pkg/scheduler"
)

func main() {
	fmt.Println("Starting Scheduler")

	// load kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", "/home/noah/.kube/config")
	if err != nil {
		fmt.Println("Error loading config")
	}

	// running in cluster
	/*
		config, err := rest.InClusterConfig()
		if err != nil {
			fmt.Println("Error creating config")
		}
	*/

	// create a kubernetes clientset from the config
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating clientset")
	}

	// create a Kubernetes informer factory
	// this creates informers
	factory := informers.NewSharedInformerFactory(
		client,
		time.Minute,
	)

	// make a pod informer with factory
	podInformer := factory.Core().V1().Pods().Informer()
	nodeInformer := factory.Core().V1().Nodes()

	// make a node Lister
	nodeLister := nodeInformer.Lister()

	// add event handlers to the informer
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// pod is currently a generic object so convert it into
			// actual Pod type
			pod := obj.(*corev1.Pod)

			if pod.Spec.NodeName == "" && pod.Spec.SchedulerName == "topsis-scheduler" {
				fmt.Println("Unscheduled pod detected")
				scheduler.SchedulePod(client, pod, nodeLister)
			}
		},
	})

	// create stop channel and start informer
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	cache.WaitForCacheSync(stopCh, podInformer.HasSynced)
	fmt.Println("Informer has started")
	<-stopCh // keeps program running indefinitely

}
