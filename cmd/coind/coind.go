package main

import (
	"log"
	"os"

	"gitee.com/iuhjui/logger"
)

func main() {
	cmdtree := NewCmdTree()
	cmdtree.SetLogger(logger.NewLogger())
	if err := cmdtree.Execute(); err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}
}
