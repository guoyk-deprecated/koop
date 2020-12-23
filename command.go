package main

import (
	"context"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	nameWildcard = "*"
)

func commandPush(ctx context.Context, cluster string, namespace string, kind string, name string) (err error) {
	var client *kubernetes.Clientset
	if _, client, err = createClient(cluster); err != nil {
		return
	}
	if _, err = client.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{}); err != nil {
		return
	}
	var kinds []string
	if kind == nameWildcard {
		kinds = knownResourceNames
	} else {
		kinds = []string{kind}
	}
	for _, kind := range kinds {
		var resource *Resource
		if resource, err = findResource(kind); err != nil {
			return
		}
		dir := filepath.Join(cluster, namespace, kind)
		var names []string
		if name == nameWildcard {
			var infos []os.FileInfo
			if infos, err = ioutil.ReadDir(dir); err != nil {
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
	}
	return
}

func commandPull(ctx context.Context, cluster string, namespace string, kind string, name string) (err error) {
	var client *kubernetes.Clientset
	if _, client, err = createClient(cluster); err != nil {
		return
	}
	if _, err = client.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{}); err != nil {
		return
	}

	var kinds []string
	if kind == nameWildcard {
		kinds = knownResourceNames
	} else {
		kinds = []string{kind}
	}
	for _, kind := range kinds {

		var resource *Resource
		if resource, err = findResource(kind); err != nil {
			return
		}

		dir := filepath.Join(cluster, namespace, kind)
		if err = os.MkdirAll(dir, 0755); err != nil {
			return
		}

		var names []string
		if name == nameWildcard {
			if names, err = resource.List(ctx, client, namespace); err != nil {
				return
			}
		} else {
			names = []string{name}
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
	}
	return
}
