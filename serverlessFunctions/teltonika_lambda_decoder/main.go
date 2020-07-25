package main

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/Ahsan-J/HexParser/model"
	"github.com/Ahsan-J/HexParser/teltonika"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response refer to https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func main() {
	lambda.Start(Handler)
}

// Handler function to handle Serverless function
func Handler(request events.APIGatewayProxyRequest) (Response, error) {

	var buf bytes.Buffer

	if request.PathParameters["hex"] == "" {
		teltonikaEmptyHex := model.Response{Status: 400, Message: "Empty hex received", Code: "THPERR001"}
		body, _ := teltonikaEmptyHex.GetMarshal(nil)
		json.HTMLEscape(&buf, body)
		return Response{
			StatusCode:      400,
			IsBase64Encoded: false,
			Body:            buf.String(),
		}, errors.New("Empty hex received")
	}

	teltonikaSuccessParse := model.Response{Status: 200, Message: "Success", Code: "THPSUC001"}
	parsedObject := teltonika.Parser(request.PathParameters["hex"])
	body, err := teltonikaSuccessParse.GetMarshal(parsedObject)

	if err != nil {
		teltonikaInvalidMarshal := model.Response{Status: 400, Message: "invalid data", Code: "THPERR001"}
		body, _ := teltonikaInvalidMarshal.GetMarshal(nil)
		json.HTMLEscape(&buf, body)
		return Response{
			StatusCode:      400,
			IsBase64Encoded: false,
			Body:            buf.String(),
		}, err
	}

	if parsedObject.Meta.IsValid == false {
		teltonikaInvalid := model.Response{Status: 400, Message: "invalid data", Code: "THPERR001"}
		body, _ := teltonikaInvalid.GetMarshal(nil)
		json.HTMLEscape(&buf, body)
		return Response{
			StatusCode:      400,
			IsBase64Encoded: false,
			Body:            buf.String(),
		}, err
	}

	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}
