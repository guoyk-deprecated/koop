package main

import (
	"context"
	"encoding/json"
	appv1 "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

func init() {
	knownResources = append(knownResources, &Resource{
		Kind: "cronjob",
		List: func(ctx context.Context, client *kubernetes.Clientset, namespace string) (names []string, err error) {
			var items *appv1.CronJobList
			if items, err = client.BatchV1beta1().CronJobs(namespace).List(ctx, metav1.ListOptions{}); err != nil {
				return
			}
			for _, item := range items.Items {
				names = append(names, item.Name)
			}
			return
		},
		GetJSON: func(ctx context.Context, client *kubernetes.Clientset, namespace, name string) (data []byte, err error) {
			var obj *appv1.CronJob
			if obj, err = client.BatchV1beta1().CronJobs(namespace).Get(ctx, name, metav1.GetOptions{}); err != nil {
				return
			}
			data, err = json.Marshal(obj)
			return
		},
		SetJSON: func(ctx context.Context, client *kubernetes.Clientset, namespace, name string, data []byte) (err error) {
			var obj appv1.CronJob
			if err = json.Unmarshal(data, &obj); err != nil {
				return
			}
			obj.Namespace = namespace
			obj.Name = name

			var current *appv1.CronJob
			if current, err = client.BatchV1beta1().CronJobs(namespace).Get(ctx, name, metav1.GetOptions{}); err != nil {
				if errors.IsNotFound(err) {
					err = nil
				} else {
					return
				}
			} else {
				if GateNoUpdate.IsOn() {
					log.Println("SKIP")
					return
				}
				obj.ResourceVersion = current.ResourceVersion
			}

			if _, err = client.BatchV1beta1().CronJobs(namespace).Update(ctx, &obj, metav1.UpdateOptions{}); err != nil {
				if errors.IsNotFound(err) {
					obj.ResourceVersion = ""
					if _, err = client.BatchV1beta1().CronJobs(namespace).Create(ctx, &obj, metav1.CreateOptions{}); err != nil {
						return
					}
				}
				return
			}
			return
		},
	})
	knownResourceNames = append(knownResourceNames, "cronjob")
}
