package httpsvr

import (
	"gitee.com/iuhjui/logger"
)

type Deps struct {
	listen string         // web服务器监听地址
	logger *logger.Logger // 日志服务
}

func NewDeps() *Deps {
	return new(Deps)
}

func (deps *Deps) SetListen(listen string) {
	deps.listen = listen
}

func (deps *Deps) SetLogger(logger *logger.Logger) {
	deps.logger = logger
}

func (deps *Deps) Verify() error {
	return nil
}
