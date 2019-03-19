package engine

import (
	"github.com/coind-io/coind/core/chain"
	"github.com/coind-io/coind/core/memtx"
)

type Deps struct {
	datadir string
	chain   *chain.Chain
	memtx   *memtx.MemTx
}

func NewDeps() *Deps {
	return new(Deps)
}

func (deps *Deps) SetDataDir(datadir string) {
	deps.datadir = datadir
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
