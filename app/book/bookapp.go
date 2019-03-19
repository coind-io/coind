package bookapp

import (
	"path"

	"github.com/CloudyKit/jet"
	"github.com/GeertJohan/go.rice"
)

type BookApp struct {
	deps  *Deps
	views *jet.Set
}

func NewCoinBook(deps *Deps) (*CoinBook, error) {
	cbook := new(CoinBook)
	cbook.deps = deps
	err := cbook.deps.Verify()
	if err != nil {
		return nil, err
	}
	err = cbook.unpackResource()
	if err != nil {
		return nil, err
	}
	cbook.bindResource()
	cbook.bindBookModule()
	return cbook, nil
}

func (cb *CoinBook) unpackResource() error {
	vbox, err := rice.FindBox("./views")
	if err != nil {
		return err
	}
	err = ExecuteUnpack(cb.deps.datadir, vbox)
	if err != nil {
		return err
	}
	sbox, err := rice.FindBox("./static")
	if err != nil {
		return err
	}
	err = ExecuteUnpack(cb.deps.datadir, sbox)
	if err != nil {
		return err
	}
	return nil
}

func (cb *CoinBook) bindResource() {
	engine := cb.deps.ginsvr.Engine()
	engine.Static("/static", path.Join(cb.deps.datadir, "./static"))
	return
}

func (cb *CoinBook) bindBookModule() {
	engine := cb.deps.ginsvr.Engine()
	book := NewBookModule(cb.deps)
	engine.GET("/", book.Status)
	engine.GET("/block/:bkhash", book.Block)
	return
}
