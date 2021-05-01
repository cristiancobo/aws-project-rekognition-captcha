package apigateway

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewJSONResponse(t *testing.T) {
	c := require.New(t)
	response := NewJSONResponse(http.StatusOK, map[string]string{"hello": "world"})

	c.Equal(http.StatusOK, response.StatusCode)
	c.Equal("application/json", response.Headers["Content-Type"])
	c.JSONEq(`{"hello":"world"}`, response.Body)
}
