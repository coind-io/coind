package coinbook

import (
	"github.com/coind-io/coind/core/chain"
	"github.com/coind-io/coind/net/gin"
)

type Deps struct {
	datadir string
	ginsvr  *gin.GinServer
	chain   *chain.Chain
}

func NewDeps() *Deps {
	return new(Deps)
}

func (deps *Deps) SetDataDir(datadir string) {
	deps.datadir = datadir
	return
}

func (deps *Deps) SetGinServer(gs *gin.GinServer) {
	deps.ginsvr = gs
	return
}

func (deps *Deps) SetChain(c *chain.Chain) {
	deps.chain = c
	return
}

func (deps *Deps) Verify() error {
	return nil
}
