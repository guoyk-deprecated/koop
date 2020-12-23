package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func createClient(cluster string) (restConfig *rest.Config, client *kubernetes.Clientset, err error) {
	var home string
	if home, err = os.UserHomeDir(); err != nil {
		return
	}
	if restConfig, err = clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".koop", "cluster-"+cluster+".yaml")); err != nil {
		return
	}
	if client, err = kubernetes.NewForConfig(restConfig); err != nil {
		return
	}
	return
}
