package logger

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const key = "logger"

// Middleware attaches a Logger instance with a request ID onto the context. It
// also logs every request along with metadata about the request.
func Middleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	l := New()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			t1 := time.Now()
			id, err := uuid.NewV4()
			if err != nil {
				return errors.WithStack(err)
			}
			log := l.ID(id.String())
			c.Set(key, log)
			if err := next(c); err != nil {
				c.Error(err)
			}
			t2 := time.Now()
			log.Root(Data{
				"status_code": c.Response().Status,
				"method":      c.Request().Method,
				"path":        c.Request().URL.Path,
				"route":       c.Path(),
				"duration":    t2.Sub(t1).Seconds() * 1000,
				"referer":     c.Request().Referer(),
				"user_agent":  c.Request().UserAgent(),
			}).Info("request handled")
			return nil
		}
	}
}

// FromContext returns a Logger from the given echo.Context. If there is no
// attached logger, then this will just return a new Logger instance.
func FromContext(c echo.Context) Logger {
	var log Logger
	log, ok := c.Get(key).(Logger)
	if !ok {
		log = New()
	}
	return log
}
