package memtx

import (
	"github.com/coind-io/coind/core/chain"
)

type Deps struct {
	chain *chain.Chain
}

func NewDeps() *Deps {
	return new(Deps)
}

func (deps *Deps) SetChain(c *chain.Chain) {
	deps.chain = c
}

func (deps *Deps) Verify() error {
	return nil
}
