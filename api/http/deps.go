package httpapi

import (
	"github.com/coind-io/coind/core/chain"
	"github.com/coind-io/coind/core/memtx"
	"github.com/coind-io/coind/net/http"
)

type Deps struct {
	httpsvr *httpsvr.HttpServer
	chain   *chain.Chain
	memtx   *memtx.MemTx
}

func NewDeps() *Deps {
	return new(Deps)
}

func (deps *Deps) SetHttpServer(hs *httpsvr.HttpServer) {
	deps.httpsvr = hs
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
