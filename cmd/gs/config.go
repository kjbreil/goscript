package main

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"os"
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
				switch r.(type) {
				case string:
					modules[k] = r.(string)
				default:
					return nil, fmt.Errorf("repo for %s not string", k)
				}

			}

		}
	}

	return modules, nil
}
