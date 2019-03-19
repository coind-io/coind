package main

import (
	"gitee.com/iuhjui/logger"
	"github.com/spf13/cobra"
)

type CmdTree struct {
	logger *logger.Logger
	root   *cobra.Command
}

func NewCmdTree() *CmdTree {
	tree := new(CmdTree)
	return tree
}

func (tree *CmdTree) makeRoot() error {
	root := new(cobra.Command)
	root.Use = "coind"
	root.Short = "Coind is a crypto currency keystone-coin of core node"
	tree.root = root
	return nil
}

func (tree *CmdTree) SetLogger(logger *logger.Logger) {
	tree.logger = logger
	return
}
