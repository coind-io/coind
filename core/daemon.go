package core

import (
	"github.com/coind-io/coind/api/http"
	"github.com/coind-io/coind/app/coinbook"
	"github.com/coind-io/coind/core/node"
	"github.com/coind-io/coind/net/gin"
)

type Daemon struct {
	deps    *Deps
	node    *node.Node
	ginsvr  *gin.GinServer
	httpapi *httpapi.HttpAPI
	bookapp *coinbook.CoinBook
}

func NewDaemon(deps *Deps) (*Daemon, error) {
	daemon := new(Daemon)
	daemon.deps = deps
	err := daemon.deps.Verify()
	if err != nil {
		return nil, err
	}
	err = daemon.makeNode()
	if err != nil {
		return nil, err
	}
	err = daemon.makeGinServer()
	if err != nil {
		return nil, err
	}
	err = daemon.makeHttpAPI()
	if err != nil {
		return nil, err
	}
	err = daemon.makeBookApp()
	if err != nil {
		return nil, err
	}
	return daemon, nil
}

func (d *Daemon) makeNode() error {
	datadir := d.deps.config.GetString("datadir")
	deps := node.NewDeps()
	deps.SetDataDir(datadir)
	n, err := node.NewNode(deps)
	if err != nil {
		return err
	}
	d.node = n
	return nil
}

func (d *Daemon) makeGinServer() error {
	listen := d.deps.config.GetString("network.http")
	deps := gin.NewDeps()
	deps.SetListen(listen)
	deps.SetLogger(d.deps.logger)
	gs, err := gin.NewGinServer(deps)
	if err != nil {
		return err
	}
	d.ginsvr = gs
	return nil
}

func (d *Daemon) makeHttpAPI() error {
	deps := httpapi.NewDeps()
	deps.SetGinServer(d.ginsvr)
	deps.SetChain(d.node.Chain())
	deps.SetMemTx(d.node.MemTx())
	ha, err := httpapi.NewHttpAPI(deps)
	if err != nil {
		return err
	}
	d.httpapi = ha
	return nil
}

func (d *Daemon) makeBookApp() error {
	datadir := d.deps.config.GetString("datadir")
	deps := coinbook.NewDeps()
	deps.SetDataDir(datadir)
	deps.SetGinServer(d.ginsvr)
	deps.SetChain(d.node.Chain())
	book, err := coinbook.NewCoinBook(deps)
	if err != nil {
		return err
	}
	d.bookapp = book
	return nil
}

func (d *Daemon) Close() error {
	err := d.ginsvr.Close()
	if err != nil {
		return err
	}
	return nil
}
