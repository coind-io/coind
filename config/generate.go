package config

import (
	"gitee.com/iuhjui/autocfg"
)

func Generate() *autocfg.Config {
	return autocfg.NewConfig(map[interface{}]interface{}{
		"datadir": "./datadir",
	})
}
