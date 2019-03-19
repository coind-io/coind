package config

import (
	"gitee.com/iuhjui/autocfg"
)

func LoadConfig(cfgname string) (*autocfg.Config, error) {
	return autocfg.LoadConfig(cfgname, &autocfg.AutoCFG{
		Generate: Generate,
		Verify:   Verify,
	})
}
