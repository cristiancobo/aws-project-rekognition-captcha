package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws-project-rekognition-captcha/shared/apigateway"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	Client rekognitioniface.RekognitionAPI
	bucket = GetString("IMAGES_BUCKET", "estudiantes-cloud-2021-3")
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
}

type Response struct {
	TextPredict string `json:"text_predict"`
}

func init() {
	client := rekognition.New(session.Must(session.NewSession()))
	xray.AWS(client.Client)

	Client = client
}

func detectText(bucket, pathName string) (string, error) {
	input := &rekognition.DetectTextInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(bucket),
				Name:   aws.String(pathName),
			},
		},
	}

	output, err := Client.DetectText(input)
	if err != nil {
		return "", err
	}

	var detectText string

	for _, text := range output.TextDetections {
		detectText += *text.DetectedText
	}

	return detectText, nil
}

func apiGatewayHandler(ctx context.Context, request *Request) (*apigateway.Response, error) {
	pathName := request.QueryStringParameters["name"]
	detectText, err := detectText(bucket, pathName)
	fmt.Printf("path %s\n", pathName)
	if err != nil {
		fmt.Printf("error to predict text %s\n", err.Error())
		return apigateway.NewErrorResponse(err), nil
	}

	Response := Response{
		TextPredict: detectText,
	}

	return apigateway.NewJSONResponse(http.StatusOK, Response), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
