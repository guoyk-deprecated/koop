package main

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"text/template"
)

const rawTemplate = `
package main

import (
	"context"
	"encoding/json"
	{{.PackageName}} {{.PackagePath|quote}}
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

func init() {
	knownResources = append(knownResources, &Resource{
		Kind: {{.Kind|quote}},
		List: func(ctx context.Context, client *kubernetes.Clientset, namespace string) (names []string, err error) {
			var items *{{.PackageName}}.{{.ListType}}
			if items, err = client.{{.Group}}().{{.Resource}}(namespace).List(ctx, metav1.ListOptions{}); err != nil {
				return
			}
			for _, item := range items.Items {
				names = append(names, item.Name)
			}
			return
		},
		GetJSON: func(ctx context.Context, client *kubernetes.Clientset, namespace, name string) (data []byte, err error) {
			var obj *{{.PackageName}}.{{.Type}}
			if obj, err = client.{{.Group}}().{{.Resource}}(namespace).Get(ctx, name, metav1.GetOptions{}); err != nil {
				return
			}
			data, err = json.Marshal(obj)
			return
		},
		SetJSON: func(ctx context.Context, client *kubernetes.Clientset, namespace, name string, data []byte) (err error) {
			if _, err = client.{{.Group}}().{{.Resource}}(namespace).Patch(ctx, name, types.StrategicMergePatchType, data, metav1.PatchOptions{}); err != nil {
				if errors.IsNotFound(err) {
					var obj {{.PackageName}}.{{.Type}}
					if err = json.Unmarshal(data, &obj); err != nil {
						return
					}
					obj.Namespace = namespace
					obj.Name = name
					if _, err = client.{{.Group}}().{{.Resource}}(namespace).Create(ctx, &obj, metav1.CreateOptions{}); err != nil {
						return
					}
				}
				return
			}
			return
		},
	})
	knownResourceNames = append(knownResourceNames, {{.Kind|quote}})
}
`

type options struct {
	Kind        string
	PackageName string
	PackagePath string
	Type        string
	ListType    string
	Group       string
	Resource    string
}

func generate(opts options) {
	var err error
	var tmpl *template.Template
	if tmpl, err = template.New("").Funcs(template.FuncMap{
		"quote": func(s string) string {
			return strconv.Quote(s)
		},
	}).Parse(rawTemplate); err != nil {
		panic(err)
	}
	buf := &bytes.Buffer{}
	if err = tmpl.Execute(buf, opts); err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile("resource_"+opts.Kind+".go", buf.Bytes(), 0644); err != nil {
		panic(err)
	}
}

func main() {
	/*
	generate(options{
		Kind:        "deployment",
		PackageName: "appv1",
		PackagePath: "k8s.io/api/apps/v1",
		Type:        "Deployment",
		ListType:    "DeploymentList",
		Group:       "AppsV1",
		Resource:    "Deployments",
	})
	generate(options{
		Kind:        "statefulset",
		PackageName: "appv1",
		PackagePath: "k8s.io/api/apps/v1",
		Type:        "StatefulSet",
		ListType:    "StatefulSetList",
		Group:       "AppsV1",
		Resource:    "StatefulSets",
	})
	generate(options{
		Kind:        "daemonset",
		PackageName: "appv1",
		PackagePath: "k8s.io/api/apps/v1",
		Type:        "DaemonSet",
		ListType:    "DaemonSetList",
		Group:       "AppsV1",
		Resource:    "DaemonSets",
	})
	generate(options{
		Kind:        "service",
		PackageName: "corev1",
		PackagePath: "k8s.io/api/core/v1",
		Type:        "Service",
		ListType:    "ServiceList",
		Group:       "CoreV1",
		Resource:    "Services",
	})
	generate(options{
		Kind:        "service",
		PackageName: "corev1",
		PackagePath: "k8s.io/api/core/v1",
		Type:        "Service",
		ListType:    "ServiceList",
		Group:       "CoreV1",
		Resource:    "Services",
	})
	generate(options{
		Kind:        "secret",
		PackageName: "corev1",
		PackagePath: "k8s.io/api/core/v1",
		Type:        "Secret",
		ListType:    "SecretList",
		Group:       "CoreV1",
		Resource:    "Secrets",
	})
	generate(options{
		Kind:        "configmap",
		PackageName: "corev1",
		PackagePath: "k8s.io/api/core/v1",
		Type:        "ConfigMap",
		ListType:    "ConfigMapList",
		Group:       "CoreV1",
		Resource:    "ConfigMaps",
	})
	generate(options{
		Kind:        "ingress",
		PackageName: "extensionsv1beta1",
		PackagePath: "k8s.io/api/extensions/v1beta1",
		Type:        "Ingress",
		ListType:    "IngressList",
		Group:       "ExtensionsV1beta1",
		Resource:    "Ingresses",
	})
	generate(options{
			Kind:        "hpa",
			PackageName: "autoscalingv2beta2",
			PackagePath: "k8s.io/api/autoscaling/v2beta2",
			Type:        "HorizontalPodAutoscaler",
			ListType:    "HorizontalPodAutoscalerList",
			Group:       "AutoscalingV2beta2",
			Resource:    "HorizontalPodAutoscalers",
	})
	 */
	generate(options{
		Kind:        "pvc",
		PackageName: "corev1",
		PackagePath: "k8s.io/api/core/v1",
		Type:        "PersistentVolumeClaim",
		ListType:    "PersistentVolumeClaimList",
		Group:       "CoreV1",
		Resource:    "PersistentVolumeClaims",
	})
}
