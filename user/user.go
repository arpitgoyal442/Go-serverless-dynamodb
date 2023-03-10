package user

import (
	"encoding/json"
	"errors"

	"github.com/arpit/go-serverless/validators"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

func FetchUser(email string, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},

		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)

	if err != nil {
		return nil, errors.New("Failed To fetch Single Record")
	}

	item := new(User)

	err = dynamodbattribute.UnmarshalMap(result.Item, item)

	if err != nil {
		return nil, errors.New("Failed To unmarshall")
	}

	return item, nil

}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]User, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	result, err := dynaClient.Scan(input)

	if err != nil {
		return nil, errors.New("Failed To Fetch all Records")
	}

	item := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)

	return item, nil

}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {

	var u User

	err := json.Unmarshal([]byte(req.Body), &u)

	if err != nil {
		return nil, errors.New("Failed To unmarshall in Create user")
	}

	// Check if email is valid
	if !validators.IsEmailValid(u.Email) {
		return nil, errors.New("Invalid Email")
	}

	//Check is user already exists
	currentUser, _ := FetchUser(u.Email, tableName, dynaClient)

	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New("User Exists")

	}

	av, err := dynamodbattribute.MarshalMap(u)

	if err != nil {
		return nil, errors.New("Couldnot Marshall in Dynamodb")
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New("Couldnot create user")
	}

	return &u,nil

}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {

	return &User{},nil

}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {

	email :=req.QueryStringParameters["email"]

	input:=&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	_,err:=dynaClient.DeleteItem(input)
	if(err!=nil){
		return errors.New("CouldNot Delete")
	}

	return nil



}
