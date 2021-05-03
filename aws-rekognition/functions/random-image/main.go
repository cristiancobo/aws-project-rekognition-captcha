package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"

	"github.com/aws-project-rekognition-captcha/shared/apigateway"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	Client s3iface.S3API
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
	Image string `json:"Image"`
	Name  string `json:"Name"`
}

func init() {
	client := s3.New(session.Must(session.NewSession()))
	xray.AWS(client.Client)

	Client = client
}

func randomImage(ctx context.Context) (string, string, error) {

	random := rand.Intn(1040)

	name := fmt.Sprintf(" %v.png", random)

	fmt.Printf("key %s\n", name)

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("/random/" + name),
	}

	output, err := Client.GetObjectWithContext(ctx, input)
	if err != nil {
		fmt.Printf("output error %s\n", err.Error())
		return "", "", err
	}

	bodyBytes, err := ioutil.ReadAll(output.Body)
	if err != nil {
		fmt.Printf("body Bytes  %s\n", err.Error())
		return "", "", err
	}

	sEnc := base64.StdEncoding.EncodeToString(bodyBytes)

	return sEnc, name, nil
}

func apiGatewayHandler(ctx context.Context) (*apigateway.Response, error) {
	Response := Response{}

	image, name, err := randomImage(ctx)
	if err != nil {
		fmt.Printf("random image error %s\n", err.Error())
		return apigateway.NewErrorResponse(apigateway.ErrInternalError), nil
	}

	Response.Name = name
	Response.Image = image

	return apigateway.NewJSONResponse(http.StatusAccepted, Response), nil

}

func main() {
	lambda.Start(apiGatewayHandler)
}
