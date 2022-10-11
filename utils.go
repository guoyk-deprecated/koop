package main

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const KeyRancherSelector = "workload.user.cattle.io/workloadselector"

var (
	int32Zero = int32(0)
)

func FixRancherWorkloadSelector(name string, sel *metav1.LabelSelector) {
	if sel == nil {
		return
	}
	if sel.MatchLabels == nil {
		sel.MatchLabels = map[string]string{}
	}
	if sel.MatchLabels[KeyRancherSelector] != "" {
		sel.MatchLabels = map[string]string{}
		sel.MatchLabels["app"] = name
	}
}

func FixRancherWorkloadPodTemplate(name string, sp *corev1.PodTemplateSpec) {
	if sp == nil {
		return
	}
	if sp.Labels == nil {
		sp.Labels = map[string]string{}
	}
	if sp.Labels[KeyRancherSelector] != "" {
		sp.Labels = map[string]string{}
		sp.Labels["app"] = name
	}
}

func FixRancherService(name string, sp *corev1.ServiceSpec) {
	if sp == nil {
		return
	}
	if sp.Selector == nil {
		return
	}
	if sp.Selector[KeyRancherSelector] != "" {
		sp.Selector = map[string]string{}
		sp.Selector["app"] = name
	}
}

func JSON2YAML(buf []byte) (out []byte, err error) {
	var m map[string]interface{}
	if err = json.Unmarshal(buf, &m); err != nil {
		return
	}
	out, err = yaml.Marshal(m)
	return
}

func YAML2JSON(buf []byte) (out []byte, err error) {
	var m map[string]interface{}
	if err = yaml.Unmarshal(buf, &m); err != nil {
		return
	}
	out, err = json.Marshal(m)
	return
}
