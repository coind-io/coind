package main

import (
	"gitee.com/iuhjui/autocfg"
	"gitee.com/iuhjui/logger"
	"github.com/spf13/cobra"

	"github.com/coind-io/coind/config"
	"github.com/coind-io/coind/version"
)

type CmdTree struct {
	logger *logger.Logger
	root   *cobra.Command
	config *autocfg.Config
}

func NewCmdTree() *CmdTree {
	tree := new(CmdTree)
	return tree
}

func (tree *CmdTree) makeRoot() error {
	root := new(cobra.Command)
	root.Use = "coind"
	root.Short = "Coind is a crypto currency keystone-coin of core node"
	root.PersistentFlags().StringP("config", "c", "", "path to an explicit configuration file")
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// 输出版本信息
		tree.logger.Infof("Coind version:  {{version}} ({{githash}})", map[string]string{
			"version": version.Version,
			"githash": version.GitHash,
		})
		// 初始配置
		cfgname, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}
		if cfgname == "" {
			cfgname = "./coind.yaml"
		}
		cfg, err := config.LoadConfig(cfgname)
		if err != nil {
			return err
		}
		tree.config = cfg
		tree.logger.Infof("using config file: {{path}}", map[string]string{
			"path": cfgname,
		})
		return nil
	}

	tree.root = root
	return nil
}

func (tree *CmdTree) SetLogger(logger *logger.Logger) {
	tree.logger = logger
	return
}
