package config

import (
	"gitee.com/iuhjui/autocfg"
)

func Generate() *autocfg.Config {
	return autocfg.NewConfig(map[interface{}]interface{}{
		"datadir": "./datadir",
		"seeds":   []interface{}{},
		"network": map[interface{}]interface{}{
			"http": "0.0.0.0:9233",
		},
	})
}
