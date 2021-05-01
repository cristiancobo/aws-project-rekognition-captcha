package apigateway

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
)

func TestParseRequestFormURLEncodedValid(t *testing.T) {
	c := require.New(t)

	values, err := ParseRequest(&events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: "a=1&b=2&c=3&c=3",
	})

	c.Nil(err)
	c.NotNil(values)
	c.Equal("1", values["a"][0])
	c.Equal("2", values["b"][0])
	c.Equal("3", values["c"][0])
}

func TestParseRequestJSONValid(t *testing.T) {
	c := require.New(t)

	values, err := ParseRequest(&events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"a": 1, "b": "2", "c": ["3", "3"]}`,
	})

	c.Nil(err)
	c.NotNil(values)
	c.Equal("1", values["a"][0])
	c.Equal("2", values["b"][0])
	c.Equal("3", values["c"][0])
}

func TestParseRequestJSONInvalid(t *testing.T) {
	c := require.New(t)

	values, err := ParseRequest(&events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `a=1`,
	})

	c.NotNil(err)
	c.Empty(values)
}
