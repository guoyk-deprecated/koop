package main

import (
	"encoding/json"
	jsonpatch "github.com/evanphx/json-patch"
)

var (
	removesDefault = []string{
		"/metadata/annotations/analysis.crane.io~1replicas-recommendation",
		"/metadata/annotations/analysis.crane.io~1resource-recommendation",
		"/metadata/annotations/net.guoyk.autodown~1lease",
		"/metadata/annotations/deployment.kubernetes.io~1revision",
		"/metadata/annotations/field.cattle.io~1creatorId",
		"/metadata/annotations/field.cattle.io~1ingressState",
		"/metadata/annotations/field.cattle.io~1publicEndpoints",
		"/metadata/annotations/field.cattle.io~1targetWorkloadIds",
		"/metadata/annotations/kubectl.kubernetes.io~1last-applied-configuration",
		"/metadata/annotations/workload.cattle.io~1targetWorkloadIdNoop",
		"/metadata/annotations/workload.cattle.io~1workloadPortBased",
		"/metadata/creationTimestamp",
		"/metadata/finalizers",
		"/metadata/namespace",
		"/metadata/generation",
		"/metadata/labels/cattle.io~1creator",
		"/metadata/labels/workload.user.cattle.io~1workloadselector",
		"/metadata/managedFields",
		"/metadata/ownerReferences",
		"/metadata/selfLink",
		"/metadata/uid",
		"/spec/replicas",
		"/spec/clusterIP",
		"/spec/clusterIPs",
		"/spec/template/metadata/annotations/cattle.io~1timestamp",
		"/spec/template/metadata/annotations/field.cattle.io~1ports",
		"/spec/template/metadata/annotations/net.guoyk.deployer~1timestamp",
		"/spec/template/metadata/annotations/workload.cattle.io~1state",
		"/spec/template/metadata/creationTimestamp",
		"/status",
	}
)

var (
	sanitizersDefault           = PatchSet{}
	sanitizersNoResourceVersion = PatchSet{
		{{Op: OpRemove, Path: "/metadata/resourceVersion"}},
	}
)

func init() {
	for _, remove := range removesDefault {
		sanitizersDefault = append(sanitizersDefault, Patches{
			Patch{Op: OpRemove, Path: remove},
		})
	}
}

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

type Patches []Patch

type PatchSet []Patches

func (ps PatchSet) Apply(data []byte) (out []byte, err error) {
	for _, item := range ps {
		var buf []byte
		if buf, err = json.Marshal(item); err != nil {
			return
		}
		var patch jsonpatch.Patch
		if patch, err = jsonpatch.DecodePatch(buf); err != nil {
			return
		}
		if buf, err = patch.Apply(data); err == nil {
			data = buf
		}
	}
	out = data
	err = nil
	return
}
