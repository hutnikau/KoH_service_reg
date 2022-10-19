package main

import (
	"encoding/json"
	"net/http"
	"service-reg/pkg/handlers"
	"service-reg/pkg/model"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	lambda.Start(handler)
}

type ErrorBody struct {
	ErrorMsg string `json:"error,omitempty"`
}

func handler(req events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	reqJson, _ := json.Marshal(req)
	log.Info().RawJSON("Raw json", reqJson).Msg("Raw json")

	switch req.RequestContext.RouteKey {
	case "POST /register":
		user := model.User{}
		_ = json.Unmarshal([]byte(req.Body), &user)

		u, err := handlers.Register(&user)
		if err != nil {
			return handlers.ApiResponse(http.StatusBadRequest, ErrorBody{err.Error()})
		}
		return handlers.ApiResponse(http.StatusAccepted, u)
	}
	return handlers.UnhandledMethod()
}
