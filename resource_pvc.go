package main

import (
	"context"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func init() {
	knownResources = append(knownResources, &Resource{
		Kind: "pvc",
		List: func(ctx context.Context, client *kubernetes.Clientset, namespace string) (names []string, err error) {
			var items *corev1.PersistentVolumeClaimList
			if items, err = client.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{}); err != nil {
				return
			}
			for _, item := range items.Items {
				names = append(names, item.Name)
			}
			return
		},
		GetJSON: func(ctx context.Context, client *kubernetes.Clientset, namespace, name string) (data []byte, err error) {
			var obj *corev1.PersistentVolumeClaim
			if obj, err = client.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, metav1.GetOptions{}); err != nil {
				return
			}
			data, err = json.Marshal(obj)
			return
		},
		SetJSON: func(ctx context.Context, client *kubernetes.Clientset, namespace, name string, data []byte) (err error) {
			var obj corev1.PersistentVolumeClaim
			if err = json.Unmarshal(data, &obj); err != nil {
				return
			}
			obj.Namespace = namespace
			obj.Name = name
			if _, err = client.CoreV1().PersistentVolumeClaims(namespace).Update(ctx, &obj, metav1.UpdateOptions{}); err != nil {
				if errors.IsNotFound(err) {
					obj.ResourceVersion = ""
					if _, err = client.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, &obj, metav1.CreateOptions{}); err != nil {
						return
					}
				}
				return
			}
			return
		},
	})
	knownResourceNames = append(knownResourceNames, "pvc")
}
