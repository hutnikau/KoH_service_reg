package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	reqJson, _ := json.Marshal(req)
	log.Debug().Msg(string(reqJson))

	log.Info().RawJSON("Raw json", reqJson).Msg("Raw json")

	return apiResponse(http.StatusAccepted, "Registered successfully")
}

func apiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}
