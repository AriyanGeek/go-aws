package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)
type Device struct {
		Id string `json:"id"`
		DeviceModel string `json:"deviceModel"`
		Name string `json:"name"`
		Note string `json:"note"`
		Serial string `json:"serial"`
	}
func Get(svc dynamodbiface.DynamoDBAPI , request events.APIGatewayProxyRequest) (*dynamodb.GetItemOutput, error) {
	type IdType struct {
		Id string `json:"id"`
	}
	var val = request.PathParameters["id"]
	id := IdType{Id:val}

	key, _ := dynamodbattribute.MarshalMap(id)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return nil, err
	// }

	param := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String("device"),
	}

	result, err := svc.GetItem(param)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func Post(svc dynamodbiface.DynamoDBAPI , request events.APIGatewayProxyRequest) (*dynamodb.PutItemOutput, error){
	
	fmt.Println("reqBody:"+request.Body)
	device := Device{}
	json.Unmarshal([]byte(request.Body), &device)
	fmt.Println(device)
	if device.Id == "" {
		return nil , errors.New("id field is empty")
	}
	fmt.Println("pass id")
	if device.DeviceModel == "" {
		return nil , errors.New("device model field is empty")
	}
	fmt.Println("pass dm")
	if device.Name == "" {
		return nil , errors.New("name field is empty")
	}
	fmt.Println("pass name")
	if device.Note == "" {
		return nil , errors.New("note field is empty")
	}
	if device.Serial == "" {
		return nil , errors.New("serial field is empty")
	}
	newItem, _ := dynamodbattribute.MarshalMap(device)
	// if err != nil {
	// 	fmt.Println(err.Error()+"line 91")
	// 	return nil, err
	// }
	input := &dynamodb.PutItemInput{
		Item:      newItem,
		TableName: aws.String("device"),
	}
	result, err := svc.PutItem(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func CreateSvc()(svc dynamodbiface.DynamoDBAPI){
	config := &aws.Config{
		Region:   aws.String("us-east-1"),
	}
	session := session.Must(session.NewSession(config))
	// Create DynamoDB client
	return dynamodb.New(session)
}
func HandleRequestProxy(svc dynamodbiface.DynamoDBAPI , request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
	
	switch request.HTTPMethod {
	case "GET":
		// call get func 
		devices := Device{}
		result, err := Get(svc , request)
		if err != nil {
			ApiResponse := events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 502}
			return ApiResponse , err
		}
		if result.Item == nil {
			ApiResponse := events.APIGatewayProxyResponse{Body: "No Item Found", StatusCode: 404}
			return ApiResponse , err
		}
		dynamodbattribute.UnmarshalMap(result.Item, &devices)
		// if err != nil {
		// 	// fmt.Println(err.Error())
		// 	ApiResponse := events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 502}
		// 	return ApiResponse, err
		// }

		responseBody , _ := json.Marshal(devices)
		// if err != nil {
		// 	// fmt.Println(err.Error())
		// 	ApiResponse := events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 502}
		// 	return ApiResponse, err
		// }
		response := events.APIGatewayProxyResponse{Body: string(responseBody), StatusCode: 200}
		return response, nil

	case "POST":
		if request.Body == "" {
			err:= errors.New("empty body")
			fmt.Println(err.Error()+"line:102")
			ApiResponse := events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}
			return ApiResponse, err
		}
		_, err := Post(svc , request)
		if err != nil {
			fmt.Println(err.Error()+"line:138")
			ApiResponse := events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}
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
func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
	svc := CreateSvc()
	return HandleRequestProxy(svc,request)
}

func main()  {
	lambda.Start(HandleRequest)
}