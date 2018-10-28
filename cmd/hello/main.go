package main

import (
	"net/http"

	"github.com/robinjoseph08/hello/pkg/application"
	"github.com/robinjoseph08/hello/pkg/logger"
	"github.com/robinjoseph08/hello/pkg/server"
)

func main() {
	log := logger.New()

	app := application.New()
	srv, err := server.New(app)
	if err != nil {
		log.Err(err).Fatal("server error")
	}

	log.Info("server started", logger.Data{"port": app.Config.Port})
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Err(err).Fatal("server stopped")
	}
	log.Info("server stopped")
}
