package apigateway

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
)

// Request is a request from the AWS API Gateway when using the default Lambda proxy.
type Request = events.APIGatewayProxyRequest

// ParseRequest parses the request body depending on the content type
func ParseRequest(req *events.APIGatewayProxyRequest) (url.Values, error) {
	if req.Headers["Content-Type"] == "application/json" {
		return parseJSONBody(req.Body)
	}

	return url.ParseQuery(req.Body)
}

func parseJSONBody(body string) (url.Values, error) {
	parsedBody := url.Values{}
	parsedJSONBody := map[string]interface{}{}

	err := json.Unmarshal([]byte(body), &parsedJSONBody)
	if err != nil {
		return parsedBody, err
	}

	for key, valueI := range parsedJSONBody {
		parsedBody[key] = interfaceAsValues(valueI)
	}

	return parsedBody, nil
}

func interfaceAsValues(i interface{}) []string {
	switch v := i.(type) {
	case string:
		{
			return []string{v}
		}
	case []interface{}:
		{
			values := make([]string, len(v))
			for i, innerValue := range v {
				values[i] = fmt.Sprintf("%v", innerValue)
			}

			return values
		}
	default:
		{
			return []string{fmt.Sprintf("%v", v)}
		}
	}
}
