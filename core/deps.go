package core

import (
	"gitee.com/iuhjui/autocfg"
	"gitee.com/iuhjui/logger"
)

type Deps struct {
	logger *logger.Logger
	config *autocfg.Config
}

func NewDeps() *Deps {
	return new(Deps)
}

func (deps *Deps) SerLogger(logger *logger.Logger) {
	deps.logger = logger
}

func (deps *Deps) SetConfig(config *autocfg.Config) {
	deps.config = config
}

func (deps *Deps) Verify() error {
	return nil
}
