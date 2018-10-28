package recovery

import (
	"fmt"

	"github.com/labstack/echo"
)

// Middleware returns recovers from any possible panics in subsequent handlers
// and funnels it to the error handler to be returned as a 500.
func Middleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err := fmt.Errorf("%v", r)
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
