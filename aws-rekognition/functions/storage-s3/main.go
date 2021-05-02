package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
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
	// Client is the client for the S3 service
	Client s3iface.S3API

	bucket       = GetString("IMAGES_BUCKET", "estudiantes-cloud-2021-3")
	mapImgFormat = map[string]format{
		"/9j/": {
			Format:      ".jpg",
			ContentType: "image/jpeg",
		},
		"iVBO": {
			Format:      ".png",
			ContentType: "image/png",
		},
		"UE5H": {
			Format:      ".png",
			ContentType: "image/png",
		},
		"R0lG": {
			Format:      ".gif",
			ContentType: "image/png",
		},
		"Qk0u": {
			Format:      ".bmp",
			ContentType: "image/bmp",
		},
		"Qk1W": {
			Format:      ".bmp",
			ContentType: "image/bmp",
		},
	}
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

type Storage struct {
	image string
	name  string
}
type format struct {
	Format      string
	ContentType string
}

type resposeAPI struct {
	Name   string `json:"Name"`
	Path   string `json:"Path"`
	Bucket string `json:"Bucket"`
}

func init() {
	client := s3.New(session.Must(session.NewSession()))
	xray.AWS(client.Client)

	Client = client
}

func uploadObject(ctx context.Context, bucket, path, contentType string, data []byte, metadata map[string]string) error {
	input := &s3.PutObjectInput{
		Body:        bytes.NewReader(data),
		Bucket:      aws.String(bucket),
		Key:         aws.String(path),
		Metadata:    aws.StringMap(metadata),
		ContentType: aws.String(contentType),
	}

	_, err := Client.PutObjectWithContext(ctx, input)
	return err
}

func getFormat(encodedBase64 string) (string, string) {
	key := encodedBase64[:4]

	format, ok := mapImgFormat[key]
	if ok {
		return format.Format, format.ContentType
	}

	return "", ""
}

func apiGatewayHandler(ctx context.Context, request *Request) (*apigateway.Response, error) {
	queryParams, err := url.ParseQuery(request.Body)
	fmt.Printf("Body %s\n", request.Body)
	if err != nil {
		return apigateway.NewErrorResponse(apigateway.ErrInvalidRequest), nil
	}

	storage := Storage{}

	storage.image = queryParams.Get("image")
	storage.name = queryParams.Get("name")

	format, contentType := getFormat(storage.image)
	fmt.Printf("format %s\n", format)
	fmt.Printf("content type %s\n", contentType)

	path := storage.name + format

	decodedImage, err := base64.StdEncoding.DecodeString(storage.image)
	if err != nil {
		fmt.Printf("decode image error %s\n", err.Error())
		return apigateway.NewErrorResponse(err), nil
	}

	err = uploadObject(ctx, bucket, path, contentType, decodedImage, nil)
	if err != nil {
		fmt.Printf("upload object error %s\n", err.Error())
		return apigateway.NewErrorResponse(err), nil
	}

	resposeAPI := resposeAPI{
		Name:   storage.name,
		Path:   path,
		Bucket: bucket,
	}

	return apigateway.NewJSONResponse(http.StatusAccepted, resposeAPI), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
