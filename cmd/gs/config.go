package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func parseConfig(filename string) (map[string]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var configMap map[string]any
	err = yaml.Unmarshal(data, &configMap)
	if err != nil {
		return nil, err
	}

	modules := make(map[string]string)

	for k, v := range configMap {
		switch vv := v.(type) {
		case map[string]any:
			if r, ok := vv["repo"]; ok {
				if rStr, ok := r.(string); ok {
					modules[k] = rStr
				} else {
					return nil, fmt.Errorf("repo for %s not string", k)
				}
			}
		}
	}

	return modules, nil
}
