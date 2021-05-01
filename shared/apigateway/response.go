package apigateway

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// Response is the same as events.APIGatewayProxyResponse, left here for compatibility purposes
type Response = events.APIGatewayProxyResponse

// NewJSONResponse creates a new JSON response given a serializable `v`
func NewJSONResponse(statusCode int, v interface{}) *events.APIGatewayProxyResponse {
	data, _ := json.Marshal(v)

	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(data),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}
