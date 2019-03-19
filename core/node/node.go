package node

import (
	"github.com/coind-io/coind/core/chain"
	"github.com/coind-io/coind/core/engine"
	"github.com/coind-io/coind/core/memtx"
)

type Node struct {
	deps   *Deps
	ticker *Ticker
	chain  *chain.Chain
	memtx  *memtx.MemTx
	engine *engine.Engine
}

func NewNode(deps *Deps) (*Node, error) {
	node := new(Node)
	node.deps = deps
	err := node.deps.Verify()
	if err != nil {
		return nil, err
	}
	err = node.makeChain()
	if err != nil {
		return nil, err
	}
	err = node.makeMemTx()
	if err != nil {
		return nil, err
	}
	err = node.makeEngine()
	if err != nil {
		return nil, err
	}
	err = node.makeTicker()
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (n *Node) makeChain() error {
	deps := chain.NewDeps()
	deps.SetDadaDir(n.deps.datadir)
	c, err := chain.NewChain(deps)
	if err != nil {
		return err
	}
	n.chain = c
	return nil
}

func (n *Node) makeMemTx() error {
	deps := memtx.NewDeps()
	deps.SetChain(n.chain)
	mtx, err := memtx.NewMemTx(deps)
	if err != nil {
		return err
	}
	n.memtx = mtx
	return nil
}

func (n *Node) makeEngine() error {
	deps := engine.NewDeps()
	deps.SetDataDir(n.deps.datadir)
	deps.SetChain(n.chain)
	deps.SetMemTx(n.memtx)
	eng, err := engine.NewEngine(deps)
	if err != nil {
		return err
	}
	n.engine = eng
	return nil
}

func (n *Node) makeTicker() error {
	n.ticker = newTicker()
	n.ticker.Bind(n.engine.MainLoop)
	return nil
}

func (n *Node) Chain() *chain.Chain {
	return n.chain
}

func (n *Node) MemTx() *memtx.MemTx {
	return n.memtx
}

func (n *Node) Close() error {
	err := n.ticker.Close()
	if err != nil {
		return err
	}
	return nil
}
