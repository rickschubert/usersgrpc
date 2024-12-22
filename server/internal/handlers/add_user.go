package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/rickschubert/usersgrpc/config"
	"github.com/rickschubert/usersgrpc/users"
	"go.uber.org/zap"
	"k8s.io/utils/ptr"
)

func AddUser(
	ctx context.Context,
	logger *zap.SugaredLogger,
	dynamoClient *dynamodb.DynamoDB,
	input *users.AddUserInput,
) (*users.User, error) {
	id := uuid.New().String()

	hashedPassword, err := hashPassword(input.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	item := map[string]*dynamodb.AttributeValue{
		"Id":        {S: aws.String(id)},
		"FirstName": {S: aws.String(input.GetFirstName())},
		"LastName":  {S: aws.String(input.GetLastName())},
		"Nickname":  {S: aws.String(input.GetNickname())},
		"Password":  {S: aws.String(string(hashedPassword))},
		"Email":     {S: aws.String(input.GetEmail())},
		"Country":   {S: aws.String(input.GetCountry())},
		"CreatedAt": {S: aws.String(time.Now().Format(time.RFC3339))},
	}

	dynamoInput := &dynamodb.PutItemInput{
		TableName: aws.String(config.TableName()),
		Item:      item,
	}

	_, err = dynamoClient.PutItem(dynamoInput)
	if err != nil {
		return nil, fmt.Errorf("failed to insert item into table: %v", err)
	}

	return &users.User{
		Id: ptr.To(id),
	}, nil
}
