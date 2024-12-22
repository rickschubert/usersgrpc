package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rickschubert/usersgrpc/config"
	"github.com/rickschubert/usersgrpc/users"
	"go.uber.org/zap"
)

func GetUser(
	ctx context.Context,
	logger *zap.SugaredLogger,
	dynamoClient *dynamodb.DynamoDB,
	opts *users.GetUserInput,
) (*users.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(config.TableName()),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(opts.GetId()),
			},
		},
	}

	result, err := dynamoClient.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	if result.Item == nil {
		logger.Warn("No user found with ID: ", opts.GetId())
		return nil, nil
	}

	return &users.User{
		Id:        result.Item["Id"].S,
		FirstName: result.Item["FirstName"].S,
		LastName:  result.Item["LastName"].S,
		Nickname:  result.Item["Nickname"].S,
		Password:  result.Item["Password"].S,
		Email:     result.Item["Email"].S,
		Country:   result.Item["Country"].S,
		CreatedAt: result.Item["CreatedAt"].S,
	}, nil
}
