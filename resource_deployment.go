package main

import (
	"context"
	"encoding/json"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func init() {
	knownResources = append(knownResources, &Resource{
		Kind: "deployment",
		List: func(ctx context.Context, client *kubernetes.Clientset, namespace string) (names []string, err error) {
			var items *appv1.DeploymentList
			if items, err = client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{}); err != nil {
				return
			}
			for _, item := range items.Items {
				names = append(names, item.Name)
			}
			return
		},
		GetJSON: func(ctx context.Context, client *kubernetes.Clientset, namespace, name string) (data []byte, err error) {
			var obj *appv1.Deployment
			if obj, err = client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{}); err != nil {
				return
			}
			data, err = json.Marshal(obj)
			return
		},
		SetJSON: func(ctx context.Context, client *kubernetes.Clientset, namespace, name string, data []byte) (err error) {
			var obj appv1.Deployment
			if err = json.Unmarshal(data, &obj); err != nil {
				return
			}
			obj.Namespace = namespace
			obj.Name = name

			if current, _ := client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{}); current != nil {
				obj.Spec.Replicas = current.Spec.Replicas
			}

			if _, err = client.AppsV1().Deployments(namespace).Update(ctx, &obj, metav1.UpdateOptions{}); err != nil {
				if errors.IsNotFound(err) {
					if _, err = client.AppsV1().Deployments(namespace).Create(ctx, &obj, metav1.CreateOptions{}); err != nil {
						return
					}
				}
				return
			}
			return
		},
	})
	knownResourceNames = append(knownResourceNames, "deployment")
}
