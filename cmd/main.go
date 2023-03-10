package  main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/arpit/go-serverless/pkg/handlers"
)

var dynaClient  dynamodbiface.DynamoDBAPI


func main(){

	region := os.Getenv("AWS_REGION")

	awsSession,err := session.NewSession(&aws.Config{
		Region: aws.String(region)})

	if err !=nil{
		return
	}

	dynaClient = dynamodb.New(awsSession)
	lambda.Start(handler)

}

const TableName= "LambdaInGoUser"


func handler(req events.APIGatewayProxyRequest) ( *events.APIGatewayProxyResponse,error){

	switch req.HTTPMethod{

	case "GET":
		return handlers.GetUser(req,TableName,dynaClient)
	case  "POST":
		return handlers.CreateUser(req,TableName,dynaClient)

	case   "PUT":
		return handlers.UpdateUser(req,TableName,dynaClient)
	case   "DELETE":
		return handlers.DeleteUser(req,TableName,dynaClient)
	default :
	    return handlers.UnhandledMethod()
	}


}