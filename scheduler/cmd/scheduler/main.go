package main

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
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
	informer := factory.Core().V1().Pods().Informer()

	// add event handlers to the informer
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// pod is currently a generic object so convert it into
			// actual Pod type
			pod := obj.(*corev1.Pod)

			if pod.Spec.NodeName == "" && pod.Spec.SchedulerName == "topsis-scheduler" {
				fmt.Println("Unscheduled pod detected")
			}
		},
	})

	// create stop channel and start informer
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	cache.WaitForCacheSync(stopCh, informer.HasSynced)
	fmt.Println("Informer has started")
	<-stopCh // keeps program running indefinitely

}
