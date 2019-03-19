package main

import (
	"gitee.com/iuhjui/logger"
)

func main() {
	cmdtree := NewCmdTree()
	cmdtree.SetLogger(logger.NewLogger())
	return
}
