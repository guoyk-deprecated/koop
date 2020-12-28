package main

import (
	"context"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	ignoredNamespaces = []string{
		"cattle-prometheus",
		"cattle-system",
		"kube-system",
		"kube-public",
		"kube-node-lease",
		"nginx-ingress",
		"ingress-nginx",
		"nfs-client-provisioner",
		"security-scan",
		"nfs-provisioner",
	}
)

const (
	wildcard     = "-"
	configDir    = ".koop"
	configPrefix = "cluster-"
	configSuffix = ".yaml"
)

func iterateCluster(cluster string, fn func(cluster string, client *kubernetes.Clientset) error) (err error) {
	var clusters []string
	if cluster == wildcard {
		var home string
		if home, err = os.UserHomeDir(); err != nil {
			return
		}
		dir := filepath.Join(home, configDir)
		var infos []os.FileInfo
		if infos, err = ioutil.ReadDir(dir); err != nil {
			return
		}
		for _, info := range infos {
			if info.IsDir() {
				continue
			}
			if !strings.HasPrefix(info.Name(), configPrefix) {
				continue
			}
			if !strings.HasSuffix(info.Name(), configSuffix) {
				continue
			}
			clusters = append(clusters, strings.TrimSuffix(strings.TrimPrefix(info.Name(), configPrefix), configSuffix))
		}
	} else {
		clusters = []string{cluster}
	}
	for _, cluster := range clusters {
		var home string
		var restConfig *rest.Config
		var client *kubernetes.Clientset
		if home, err = os.UserHomeDir(); err != nil {
			return
		}
		if restConfig, err = clientcmd.BuildConfigFromFlags("", filepath.Join(home, configDir, configPrefix+cluster+configSuffix)); err != nil {
			return
		}
		if client, err = kubernetes.NewForConfig(restConfig); err != nil {
			return
		}
		if err = fn(cluster, client); err != nil {
			return
		}
	}
	return
}

func iterateNamespace(ctx context.Context, client *kubernetes.Clientset, namespace string, fn func(namespace string) error) (err error) {
	var namespaces []string
	if namespace == wildcard {
		var items *corev1.NamespaceList
		if items, err = client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{}); err != nil {
			return
		}
	outerLoop:
		for _, item := range items.Items {
			for _, ignored := range ignoredNamespaces {
				if strings.HasPrefix(strings.ToLower(item.Name), ignored) {
					continue outerLoop
				}
			}
			namespaces = append(namespaces, item.Name)
		}
	} else {
		if _, err = client.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{}); err != nil {
			return
		}
		namespaces = []string{namespace}
	}
	for _, namespace := range namespaces {
		if err = fn(namespace); err != nil {
			return
		}
	}
	return
}

func iterateKind(kind string, fn func(kind string) error) (err error) {
	var kinds []string
	if kind == wildcard {
		kinds = knownResourceNames
	} else {
		kinds = []string{kind}
	}
	for _, kind := range kinds {
		if err = fn(kind); err != nil {
			return
		}
	}
	return
}

func commandPush(ctx context.Context, cluster string, namespace string, kind string, name string) (err error) {
	if err = iterateCluster(cluster, func(cluster string, client *kubernetes.Clientset) error {
		return iterateNamespace(ctx, client, namespace, func(namespace string) error {
			return iterateKind(kind, func(kind string) (err error) {
				var resource *Resource
				if resource, err = findResource(kind); err != nil {
					return
				}
				dir := filepath.Join(cluster, namespace, kind)
				var names []string
				if name == wildcard {
					var infos []os.FileInfo
					if infos, err = ioutil.ReadDir(dir); err != nil {
						if os.IsNotExist(err) {
							err = nil
						}
						return
					}
					for _, info := range infos {
						if info.IsDir() {
							log.Println("found unexpected directory in:", dir)
							continue
						}
						if !strings.HasSuffix(info.Name(), ".yaml") {
							log.Println("found unexpected file", info.Name(), "in:", dir, ", for compatible reasons, all YAML files must has extension '.yaml', NOT '.yml'")
							continue
						}
						names = append(names, strings.TrimSuffix(info.Name(), ".yaml"))
					}
				} else {
					names = []string{name}
				}
				for _, name := range names {
					log.Printf("PUSH: %s/%s/%s/%s", cluster, namespace, kind, name)
					var buf []byte
					if buf, err = ioutil.ReadFile(filepath.Join(dir, name+".yaml")); err != nil {
						return
					}
					if err = resource.SetCanonicalYAML(ctx, client, namespace, name, buf); err != nil {
						return
					}
				}
				return
			})
		})
	}); err != nil {
		return
	}
	return
}

func commandPull(ctx context.Context, cluster string, namespace string, kind string, name string) (err error) {
	if err = iterateCluster(cluster, func(cluster string, client *kubernetes.Clientset) error {
		return iterateNamespace(ctx, client, namespace, func(namespace string) error {
			return iterateKind(kind, func(kind string) (err error) {
				var resource *Resource
				if resource, err = findResource(kind); err != nil {
					return
				}

				dir := filepath.Join(cluster, namespace, kind)

				var names []string
				if name == wildcard {
					_ = os.RemoveAll(dir)
					log.Printf("CLEAN: %s/%s/%s", cluster, namespace, kind)
					if names, err = resource.List(ctx, client, namespace); err != nil {
						return
					}
				} else {
					names = []string{name}
				}

				if err = os.MkdirAll(dir, 0755); err != nil {
					return
				}

				for _, name := range names {
					log.Printf("PULL: %s/%s/%s/%s", cluster, namespace, kind, name)
					var buf []byte
					if buf, err = resource.GetCanonicalYAML(ctx, client, namespace, name); err != nil {
						return
					}
					if err = ioutil.WriteFile(filepath.Join(dir, name+".yaml"), buf, 0755); err != nil {
						return
					}
				}
				return
			})
		})
	}); err != nil {
		return
	}
	return
}
