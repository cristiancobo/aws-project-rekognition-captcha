package main

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/aws-project-rekognition-captcha/shared/apigateway"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	awsS3ARN = GetString("AWS_S3_ARN", "arn:aws:s3:::estudiantes-cloud-2021-3")
)

// GetString gets the env var as a string
func GetString(varName string, defaultValue string) string {
	val, _ := os.LookupEnv(varName)
	if val == "" {
		return defaultValue
	}

	return val
}

type Request struct {
	*events.APIGatewayProxyRequest
	startingTime time.Time
	err          error
}

type Storage struct {
	image string
	name  string
}

func apiGatewayHandler(ctx context.Context, request *Request) (*apigateway.Response, error) {
	queryParams, err := url.ParseQuery(request.Body)
	if err != nil {
		return apigateway.NewErrorResponse(apigateway.ErrInvalidRequest), nil
	}

	storage := Storage{}

	storage.image = queryParams.Get("image")
	storage.name = queryParams.Get("name")

	// TO DO, save in S3.

	return nil, nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
