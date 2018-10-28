package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	elog "github.com/labstack/gommon/log"
	"github.com/robinjoseph08/hello/pkg/application"
	"github.com/robinjoseph08/hello/pkg/health"
	"github.com/robinjoseph08/hello/pkg/logger"
	"github.com/robinjoseph08/hello/pkg/recovery"
	"github.com/robinjoseph08/hello/pkg/signals"
)

// New returns a new HTTP server with the registered routes.
func New(app application.App) (*http.Server, error) {
	log := logger.New()
	e := echo.New()

	e.Logger.SetLevel(elog.OFF)

	e.Use(logger.Middleware())
	e.Use(recovery.Middleware())

	health.RegisterRoutes(e)

	echo.NotFoundHandler = notFoundHandler

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Port),
		Handler: e,
	}

	graceful := signals.Setup()

	go func() {
		<-graceful
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Err(err).Error("server shutdown")
		}
	}()

	return srv, nil
}

func notFoundHandler(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotFound, "not found")
}
