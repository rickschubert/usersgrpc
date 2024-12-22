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

func RemoveUser(
	ctx context.Context,
	logger *zap.SugaredLogger,
	dynamoClient *dynamodb.DynamoDB,
	opts *users.RemoveUserInput,
) (*users.User, error) {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(config.TableName()),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(opts.GetId()),
			},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueAllOld),
	}

	result, err := dynamoClient.DeleteItem(input)
	if err != nil {
		return nil, fmt.Errorf("failed to delete item: %v", err)
	}

	if result.Attributes["Id"] != nil {
		return &users.User{
			Id:        result.Attributes["Id"].S,
			FirstName: result.Attributes["FirstName"].S,
			LastName:  result.Attributes["LastName"].S,
			Nickname:  result.Attributes["Nickname"].S,
			Password:  result.Attributes["Password"].S,
			Email:     result.Attributes["Email"].S,
			Country:   result.Attributes["Country"].S,
			CreatedAt: result.Attributes["CreatedAt"].S,
		}, nil
	} else {
		return nil, nil
	}
}
