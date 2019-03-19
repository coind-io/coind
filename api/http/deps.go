package httpapi

import (
	"github.com/coind-io/coind/core/chain"
	"github.com/coind-io/coind/core/memtx"
	"github.com/coind-io/coind/net/gin"
)

type Deps struct {
	ginsvr *gin.GinServer
	chain  *chain.Chain
	memtx  *memtx.MemTx
}

func NewDeps() *Deps {
	return new(Deps)
}

func (deps *Deps) SetGinServer(gs *gin.GinServer) {
	deps.ginsvr = gs
	return
}

func (deps *Deps) SetChain(c *chain.Chain) {
	deps.chain = c
}

func (deps *Deps) SetMemTx(mtx *memtx.MemTx) {
	deps.memtx = mtx
}

func (deps *Deps) Verify() error {
	return nil
}
