package httpapi

type HttpAPI struct {
	deps *Deps
}

func NewHttpAPI(deps *Deps) (*HttpAPI, error) {
	ha := new(HttpAPI)
	ha.deps = deps
	err := ha.deps.Verify()
	if err != nil {
		return nil, err
	}
	ha.bindFetchModule()
	//ha.bindSendModule()
	return ha, nil
}

func (ha *HttpAPI) bindFetchModule() {
	router := ha.deps.httpsvr.Router()
	fetch := NewFetchModule(ha.deps)
	router.Get("/api/v1/fetch/balance/{owner}").HandlerFunc(fetch.Balance)
	//router.GET("/api/v1/fetch/balance/:owner", fetch.Balance)
	//router.GET("/api/v1/fetch/coins/:owner", fetch.Coins)
	/*
		router.GetFunc("/api/v1/fetch/blocks", fetch.Blocks)
		router.GetFunc("/api/v1/fetch/block-index/:index", fetch.BlockIndex)
		router.GetFunc("/api/v1/fetch/block/:hash", fetch.Block)
		router.GetFunc("/api/v1/fetch/tx/:hash", fetch.Tx)
	*/
	return
}

/*
func (ha *HttpAPI) bindSendModule() {
	engine := ha.deps.ginsvr.Engine()
	send := NewSendModule(ha.deps)
	engine.POST("/api/v1/send/rawtx", send.RawTx)
	return
}
*/
