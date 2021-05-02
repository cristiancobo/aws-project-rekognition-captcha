package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApiGatewayHandler(t *testing.T) {
	c := require.New(t)

	input := &Request{}

	inputJSON, err := ioutil.ReadFile("samples/input-apigw.json")
	c.Nil(err)

	err = json.Unmarshal(inputJSON, &input)
	c.Nil(err)

	response, err := apiGatewayHandler(context.Background(), input)
	c.Equal(http.StatusOK, response.StatusCode)
	c.Nil(err)
}
