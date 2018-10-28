package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/robinjoseph08/hello/pkg/application"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	app := application.New()
	srv, err := New(app)
	require.NoError(t, err)

	req, err := http.NewRequest("GET", "/health", nil)
	require.Nil(t, err, "unexpecetd error when making new request")

	w := httptest.NewRecorder()

	srv.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "incorrect status code")
	assert.Equal(t, `{"healthy":true}`, w.Body.String(), "incorrect response")

	req, err = http.NewRequest("GET", "/foo", nil)
	require.Nil(t, err, "unexpecetd error when making new request")

	w = httptest.NewRecorder()

	srv.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "incorrect status code")
	assert.Contains(t, w.Body.String(), "not found", "incorrect response")
}
