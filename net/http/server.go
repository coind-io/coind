package httpsvr

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type HttpServer struct {
	deps   *Deps
	server *http.Server
	router *mux.Router
}

func NewHttpServer(deps *Deps) (*HttpServer, error) {
	hs := new(HttpServer)
	hs.deps = deps
	err := hs.deps.Verify()
	if err != nil {
		return nil, err
	}
	go hs.serve()
	return hs, nil
}

func (hs *HttpServer) serve() {
	logger := hs.deps.logger
	logger.Infof("listening on {{laddr}}, web interface at http://{{laddr}}", map[string]string{
		"laddr": hs.deps.listen,
	})
	hs.router = mux.NewRouter()
	hs.server = &http.Server{
		Addr:    hs.deps.listen,
		Handler: hs.router,
	}
	err := hs.server.ListenAndServe()
	if err != nil {
		logger.Infof("listening on {{laddr}} fail, {{error}}", map[string]string{
			"laddr": hs.deps.listen,
			"error": err.Error(),
		})
	}
	return
}

func (hs *HttpServer) exit() error {
	logger := hs.deps.logger
	logger.Infof("waiting for the remaining connections to finish...", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := hs.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	logger.Infof("gracefully shutdown the http server...", nil)
	return nil
}

func (hs *HttpServer) Router() *mux.Router {
	return hs.router
}

func (hs *HttpServer) Close() error {
	err := hs.exit()
	if err != nil {
		return err
	}
	return nil
}
