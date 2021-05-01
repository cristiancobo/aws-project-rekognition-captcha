package main

import (
	"context"
	"time"

	"github.com/aws-project-rekognition-captcha/shared/apigateway"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	*events.APIGatewayProxyRequest
	startingTime time.Time
	err          error
}

func apiGatewayHandler(ctx context.Context, request *Request) (*apigateway.Response, error) {
	return nil, nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
