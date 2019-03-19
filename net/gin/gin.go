package gin

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GinServer struct {
	deps   *Deps
	server *http.Server
	engine *gin.Engine
}

func NewGinServer(deps *Deps) (*GinServer, error) {
	gs := new(GinServer)
	gs.deps = deps
	gs.engine = gin.New()
	err := gs.deps.Verify()
	if err != nil {
		return nil, err
	}
	go gs.serve()
	return gs, nil
}

func (gs *GinServer) serve() {
	logger := gs.deps.logger
	logger.Infof("listening on {{laddr}}, web interface at http://{{laddr}}", map[string]string{
		"laddr": gs.deps.listen,
	})
	gs.server = &http.Server{
		Addr:    gs.deps.listen,
		Handler: gs.engine,
	}
	err := gs.server.ListenAndServe()
	if err != nil {
		logger.Infof("listening on {{laddr}} fail, {{error}}", map[string]string{
			"laddr": gs.deps.listen,
			"error": err.Error(),
		})
	}
	return
}

func (gs *GinServer) exit() error {
	logger := gs.deps.logger
	logger.Infof("waiting for the remaining connections to finish...", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := gs.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	logger.Infof("gracefully shutdown the http server...", nil)
	return nil
}

func (gs *GinServer) Engine() *gin.Engine {
	return gs.engine
}

func (gs *GinServer) Close() error {
	err := gs.exit()
	if err != nil {
		return err
	}
	return nil
}
