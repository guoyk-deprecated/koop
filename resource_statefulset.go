package main

import (
	"context"
	"encoding/json"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

func init() {
	knownResources = append(knownResources, &Resource{
		Kind: "statefulset",
		List: func(ctx context.Context, client *kubernetes.Clientset, namespace string) (names []string, err error) {
			var items *appv1.StatefulSetList
			if items, err = client.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{}); err != nil {
				return
			}
			for _, item := range items.Items {
				names = append(names, item.Name)
			}
			return
		},
		GetJSON: func(ctx context.Context, client *kubernetes.Clientset, namespace, name string) (data []byte, err error) {
			var obj *appv1.StatefulSet
			if obj, err = client.AppsV1().StatefulSets(namespace).Get(ctx, name, metav1.GetOptions{}); err != nil {
				return
			}
			data, err = json.Marshal(obj)
			return
		},
		SetJSON: func(ctx context.Context, client *kubernetes.Clientset, namespace, name string, data []byte) (err error) {
			if _, err = client.AppsV1().StatefulSets(namespace).Patch(ctx, name, types.StrategicMergePatchType, data, metav1.PatchOptions{}); err != nil {
				if errors.IsNotFound(err) {
					var obj appv1.StatefulSet
					if err = json.Unmarshal(data, &obj); err != nil {
						return
					}
					obj.Namespace = namespace
					obj.Name = name
					if _, err = client.AppsV1().StatefulSets(namespace).Create(ctx, &obj, metav1.CreateOptions{}); err != nil {
						return
					}
				}
				return
			}
			return
		},
	})
	knownResourceNames = append(knownResourceNames, "statefulset")
}
