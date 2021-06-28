package main

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

var (
	int32Zero = int32(0)
)

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

func IsEnvNoUpdate() bool {
	v, _ := strconv.ParseBool(os.Getenv("KOOP_NO_UPDATE"))
	return v
}

func IsEnvZeroReplicas() bool {
	v, _ := strconv.ParseBool(os.Getenv("KOOP_ZERO_REPLICAS"))
	return v
}
