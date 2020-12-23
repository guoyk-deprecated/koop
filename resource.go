package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/k8s-autoops/koop/pkg/jsonpatch"
	"k8s.io/client-go/kubernetes"
	"strings"
)

const (
	OpAdd     = "add"
	OpRemove  = "remove"
	OpReplace = "replace"
	OpCopy    = "copy"
	OpMove    = "move"
	OpTest    = "test"
)

type Patch struct {
	Op    string      `json:"op,omitempty"`
	Path  string      `json:"path,omitempty"`
	From  string      `json:"from,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

var (
	knownPatches = []Patch{
		{Op: OpRemove, Path: "/status"},
		{Op: OpRemove, Path: "/kind"},
		{Op: OpRemove, Path: "/apiVersion"},
		{Op: OpRemove, Path: "/metadata/name"},
		{Op: OpRemove, Path: "/metadata/namespace"},
		{Op: OpRemove, Path: "/metadata/creationTimestamp"},
		{Op: OpRemove, Path: "/metadata/generation"},
		{Op: OpRemove, Path: "/metadata/resourceVersion"},
		{Op: OpRemove, Path: "/metadata/selfLink"},
		{Op: OpRemove, Path: "/metadata/uid"},
		{Op: OpRemove, Path: "/metadata/annotations/kubectl.kubernetes.io~1last-applied-configuration"},
		{Op: OpRemove, Path: "/metadata/annotations/deployment.kubernetes.io~1revision"},
		{Op: OpRemove, Path: "/metadata/annotations/field.cattle.io~1ingressState"},
		{Op: OpRemove, Path: "/metadata/annotations/field.cattle.io~publicEndpoints"},
		{Op: OpRemove, Path: "/spec/template/metadata/creationTimestamp"},
		{Op: OpRemove, Path: "/spec/replicas"},
	}
)

type Resource struct {
	Kind    string
	List    func(ctx context.Context, client *kubernetes.Clientset, namespace string) ([]string, error)
	GetJSON func(ctx context.Context, client *kubernetes.Clientset, namespace, name string) ([]byte, error)
	SetJSON func(ctx context.Context, client *kubernetes.Clientset, namespace, name string, data []byte) error
}

func (r Resource) GetCanonicalYAML(ctx context.Context, client *kubernetes.Clientset, namespace, name string) (data []byte, err error) {
	if data, err = r.GetJSON(ctx, client, namespace, name); err != nil {
		return
	}
	var buf []byte
	if buf, err = json.Marshal(knownPatches); err != nil {
		return
	}
	var patch jsonpatch.Patch
	if patch, err = jsonpatch.DecodePatch(buf); err != nil {
		return
	}
	if data, err = patch.Apply(data); err != nil {
		return
	}
	data, err = JSON2YAML(data)
	return
}

func (r Resource) SetCanonicalYAML(ctx context.Context, client *kubernetes.Clientset, namespace, name string, data []byte) (err error) {
	var buf []byte
	if buf, err = json.Marshal(knownPatches); err != nil {
		return
	}
	var patch jsonpatch.Patch
	if patch, err = jsonpatch.DecodePatch(buf); err != nil {
		return
	}
	if data, err = YAML2JSON(data); err != nil {
		return
	}
	if data, err = patch.Apply(data); err != nil {
		return
	}
	if err = r.SetJSON(ctx, client, namespace, name, data); err != nil {
		return
	}
	return
}

var (
	knownResources     []*Resource
	knownResourceNames []string
)

func findResource(kind string) (resource *Resource, err error) {
	for _, knownResource := range knownResources {
		if knownResource.Kind == kind {
			resource = knownResource
		}
	}
	if resource == nil {
		err = fmt.Errorf("unknown resource kind '%s', known kinds are %s", kind, strings.Join(knownResourceNames, ", "))
		return
	}
	return
}
