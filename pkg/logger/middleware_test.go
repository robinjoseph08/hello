package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/robinjoseph08/hello/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	e := echo.New()
	e.Use(Middleware())

	e.GET("/", func(c echo.Context) error {
		log, ok := c.Get("logger").(Logger)
		assert.True(t, ok, "expected logger to be of type Logger")
		assert.NotNil(t, log)
		return nil
	})

	req, err := http.NewRequest("GET", "/", nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()

	e.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestFromContext(t *testing.T) {
	log := New().ID("foo")
	c, _ := test.NewContext(t, nil)
	c.Set(key, log)

	l := FromContext(c)

	assert.Equal(t, log.id, l.id)

	c, _ = test.NewContext(t, nil)

	l = FromContext(c)

	assert.Equal(t, "", l.id)
}
