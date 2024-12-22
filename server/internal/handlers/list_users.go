package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rickschubert/usersgrpc/config"
	"github.com/rickschubert/usersgrpc/users"
	"go.uber.org/zap"
	"k8s.io/utils/ptr"
)

func ListUsers(
	ctx context.Context,
	logger *zap.SugaredLogger,
	dynamoClient *dynamodb.DynamoDB,
	opts *users.ListUsersInput,
) (*users.ListUsersResponse, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(config.TableName()),
		Limit:     aws.Int64(int64(config.PaginationSize())),
	}

	if opts.NextPageToken != nil {
		input.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"Id": {
				S: ptr.To(opts.GetNextPageToken()),
			},
		}
	}

	if opts.Country != nil {
		input.FilterExpression = aws.String("Country = :country")
		input.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
			":country": {
				S: aws.String(*opts.Country),
			},
		}
	}

	result, err := dynamoClient.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan table: %w", err)
	}

	usersList := make([]*users.User, 0, config.PaginationSize())

	logger.Infow("Found users", zap.Int("number_of_users_found", len(result.Items)))
	for _, item := range result.Items {
		usersList = append(usersList, &users.User{
			Id:        item["Id"].S,
			FirstName: item["FirstName"].S,
			LastName:  item["LastName"].S,
			Nickname:  item["Nickname"].S,
			Password:  item["Password"].S,
			Email:     item["Email"].S,
			Country:   item["Country"].S,
			CreatedAt: item["CreatedAt"].S,
		})
	}

	var nextPageToken *string

	if len(usersList) == config.PaginationSize() {
		nextPageToken = result.LastEvaluatedKey["Id"].S
	}

	return &users.ListUsersResponse{
		Users:         usersList,
		NextPageToken: nextPageToken,
	}, nil
}
