package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
)

// Define stub
type stubDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func TestHandler(t *testing.T) {
	id1 := map[string] string{
		"id" : "id1",
	}
	id2 := map[string] string{
		"id" : "id5",
	}
	id3 := map[string] string{
		"id" : "idX",
	}
	device1 := Device{
		Id : "id100",
		DeviceModel: "cc",
		Name: "nnn",
		Note: "bbb",
		Serial: "sdsdfsfsd",
	}
	device2 := Device{
		Id : "",
		DeviceModel: "cc",
		Name: "nnn",
		Note: "bbb",
		Serial: "sdsdfsfsd",
	}
	device3 := Device{
		Id : "id100",
		DeviceModel: "",
		Name: "nnn",
		Note: "bbb",
		Serial: "sdsdfsfsd",
	}
	device4 := Device{
		Id : "id100",
		DeviceModel: "xxx",
		Name: "",
		Note: "bbb",
		Serial: "sdsdfsfsd",
	}
	device5 := Device{
		Id : "id100",
		DeviceModel: "xxx",
		Name: "nnn",
		Note: "",
		Serial: "sdsdfsfsd",
	}
	device6 := Device{
		Id : "id100",
		DeviceModel: "xxx",
		Name: "nnn",
		Note: "dsss",
		Serial: "",
	}
	device7 := Device{
		Id : "id1",
		DeviceModel: "xxx",
		Name: "nnn",
		Note: "dsss",
		Serial: "sss",
	}
	
	body1, _ := json.Marshal(device1)
	body2, _ := json.Marshal(device2)
	body3, _ := json.Marshal(device3)
	body4, _ := json.Marshal(device4)
	body5, _ := json.Marshal(device5)
	body6, _ := json.Marshal(device6)
	body7, _ := json.Marshal(device7)
	// id4 := map[string] string{
	// 	"id" : "idMarshal",
	// }
	svc := &stubDynamoDB{}
	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  int
		err     error
	}{
		{
			request: events.APIGatewayProxyRequest{HTTPMethod: "GET",PathParameters: id1},
			expect:   200,
			err:     nil,
		},
		{
			request: events.APIGatewayProxyRequest{HTTPMethod: "GET",PathParameters: id2},
			expect:   404,
			err:     nil,
		},
		{
			request: events.APIGatewayProxyRequest{HTTPMethod: "GET",PathParameters: id3},
			expect:   502,
			err:     errors.New("502"),
		},
		{
			
			request: events.APIGatewayProxyRequest{HTTPMethod: "POST"},
			expect:  400,
			err:     errors.New("empty body"),
		},
		{
			
			request: events.APIGatewayProxyRequest{HTTPMethod: "POST",Body: string(body1)},
			expect:  201,
			err:     nil,
		},
		{
			
			request: events.APIGatewayProxyRequest{HTTPMethod: "POST",Body: string(body2)},
			expect:  400,
			err:     errors.New(""),
		},
		{
			
			request: events.APIGatewayProxyRequest{HTTPMethod: "POST",Body: string(body3)},
			expect:  400,
			err:     errors.New(""),
		},
		{
			
			request: events.APIGatewayProxyRequest{HTTPMethod: "POST",Body: string(body4)},
			expect:  400,
			err:     errors.New(""),
		},
		{
			
			request: events.APIGatewayProxyRequest{HTTPMethod: "POST",Body: string(body5)},
			expect:  400,
			err:     errors.New(""),
		},
		{
			
			request: events.APIGatewayProxyRequest{HTTPMethod: "POST",Body: string(body6)},
			expect:  400,
			err:     errors.New(""),
		},
		{
			
			request: events.APIGatewayProxyRequest{HTTPMethod: "POST",Body: string(body7)},
			expect:  400,
			err:     errors.New(""),
		},
		{
			
			request: events.APIGatewayProxyRequest{HTTPMethod: "UPDATE"},
			expect:  502,
			err:     errors.New("method not allowed"),
		},
	}

	for _, test := range tests {
		response, err := HandleRequestProxy(svc ,test.request)
		assert.IsType(t, test.err, err)
		res := response.StatusCode
		assert.Equal(t, test.expect, res)
	}
}

func (m *stubDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	// Make response
	statusCode := dynamodb.AttributeValue{}
	resp := make(map[string]*dynamodb.AttributeValue)

	id := *input.Key["id"].S
	if id == "id1" {
		statusCode.SetS("200")
		resp["statusCode"] = &statusCode

		output := &dynamodb.GetItemOutput{
		Item: resp,
		}
		return output, nil
	}else if id == "idX"{
		statusCode.SetS("500")
		resp["statusCode"] = &statusCode

		return nil, errors.New("500")
	
	}else{
		statusCode.SetS("404")
		resp["statusCode"] = &statusCode
		output := &dynamodb.GetItemOutput{
		Item: nil,
		}
		return output, nil
	}
}
func (m *stubDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	// Make response
	statusCode := dynamodb.AttributeValue{}
	
	resp := make(map[string]*dynamodb.AttributeValue)

	item := input.Item
	id :=*item["id"].S
	if id == "id100" {
		statusCode.SetS("201")
		resp["statusCode"] = &statusCode
		output := &dynamodb.PutItemOutput{
			Attributes: resp,
		}
		return output, nil
	}else{
		return nil, errors.New("400")
	}
}
func (m *stubDynamoDB) HandleRequestProxy(svc dynamodbiface.DynamoDBAPI , request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
	
	switch request.HTTPMethod {
	case "GET":
		// call get func 
		devices := Device{}
		result, err := Get(svc , request)
		if err != nil {
			fmt.Println(err.Error()+":line86Test")
			ApiResponse := events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 502}
			return ApiResponse , err
		}
		if result.Item == nil {
			fmt.Println("no item found")
			ApiResponse := events.APIGatewayProxyResponse{Body: "No Item Found", StatusCode: 404}
			return ApiResponse , err
		}
		err = dynamodbattribute.UnmarshalMap(result.Item, &devices)
		
		if err != nil {
			fmt.Println(err.Error())
			ApiResponse := events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 502}
			return ApiResponse, err
		}

		responseBody , err := json.Marshal(devices)
		if err != nil {
			fmt.Println(err.Error())
			ApiResponse := events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 502}
			return ApiResponse, err
		}
		response := events.APIGatewayProxyResponse{Body: string(responseBody), StatusCode: 200}
		return response, nil

	case "POST":
		if request.Body == "" {
			err:= errors.New("Empty body")
			fmt.Println(err.Error()+"line:102")
			ApiResponse := events.APIGatewayProxyResponse{Body: err.Error()+"line:222T", StatusCode: 400}
			return ApiResponse, err
		}
		_, err := Post(svc , request)
		if err != nil {
			fmt.Println(err.Error()+"line:102")
			ApiResponse := events.APIGatewayProxyResponse{Body: err.Error()+"line:103", StatusCode: 400}
			return ApiResponse, err
		}

		fmt.Printf("We have inserted a new item!\n")
		response := events.APIGatewayProxyResponse{Body: "We have inserted a new item!", StatusCode: 201}
		return response, nil

	default:
		err := errors.New("method not allowed")
		ApiResponse := events.APIGatewayProxyResponse{Body: "Method Not ACCEPTED", StatusCode: 502}
		return ApiResponse, err

	}
}