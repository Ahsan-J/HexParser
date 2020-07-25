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

// BodyRequest will be used to take the json response from client and build it
type BodyRequest struct {
	CodecID string `json:"codecId"`
	Command string `json:"command"`
	IMEI    int    `json:"imei"`
}

// BodyResponse will handle the response structure of

func main() {
	lambda.Start(Handler)
}

// Handler function to handle Serverless function
func Handler(request events.APIGatewayProxyRequest) (Response, error) {

	var buf bytes.Buffer
	bodyRequest := BodyRequest{}

	// Unmarshal the json, return 404 if error
	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		requestBodyError := model.Response{Status: 400, Message: err.Error(), Code: "THGERR001"}
		body, _ := requestBodyError.GetMarshal(nil)
		json.HTMLEscape(&buf, body)
		return Response{
			Body:            buf.String(),
			IsBase64Encoded: false,
			StatusCode:      400,
		}, err
	}

	if bodyRequest.CodecID == "" {
		emptyCodec := model.Response{Status: 400, Message: "Need CodecID to specify the encoded string version", Code: "THPERR001"}
		body, _ := emptyCodec.GetMarshal(nil)
		json.HTMLEscape(&buf, body)
		return Response{
			Body:            buf.String(),
			IsBase64Encoded: false,
			StatusCode:      400,
		}, errors.New("Need CodecID to specify the encoded string version")
	}

	if bodyRequest.Command == "" {
		emptyCodec := model.Response{Status: 400, Message: "Invalid Command", Code: "THPERR001"}
		body, _ := emptyCodec.GetMarshal(nil)
		json.HTMLEscape(&buf, body)
		return Response{
			Body:            buf.String(),
			IsBase64Encoded: false,
			StatusCode:      400,
		}, errors.New("Invalid Command")
	}

	if bodyRequest.IMEI == 0 && bodyRequest.CodecID == "0E" {
		emptyCodec := model.Response{Status: 400, Message: "Need Device IMEI to encod in codec14", Code: "THPERR001"}
		body, _ := emptyCodec.GetMarshal(nil)
		json.HTMLEscape(&buf, body)
		return Response{
			Body:            buf.String(),
			IsBase64Encoded: false,
			StatusCode:      400,
		}, errors.New("Need Device IMEI to encod in codec14")
	}

	teltonikaSuccessEncode := model.Response{Status: 200, Message: "Success", Code: "THPSUC001"}
	parsedObject := teltonika.GenerateHex(bodyRequest.Command, bodyRequest.CodecID, bodyRequest.IMEI)
	body, err := teltonikaSuccessEncode.GetMarshal(parsedObject)

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
