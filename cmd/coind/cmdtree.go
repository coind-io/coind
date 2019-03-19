package main

import (
	"gitee.com/iuhjui/logger"
)

type CmdTree struct {
	logger *logger.Logger
}

func NewCmdTree() *CmdTree {
	tree := new(CmdTree)
	return tree
}

func (tree *CmdTree) SetLogger(logger *logger.Logger) {
	tree.logger = logger
	return
}
