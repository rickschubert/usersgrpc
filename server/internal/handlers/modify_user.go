package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rickschubert/usersgrpc/config"
	"github.com/rickschubert/usersgrpc/users"
	"go.uber.org/zap"
)

func ModifyUser(
	ctx context.Context,
	logger *zap.SugaredLogger,
	dynamoClient *dynamodb.DynamoDB,
	opts *users.ModifyUserInput,
) (*users.User, error) {
	expressionAttributeValues, err := generateExpressionAttributeValues(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate expression attribute values: %w", err)
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(config.TableName()),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(opts.GetId()),
			},
		},
		UpdateExpression:          aws.String(generateUpdateExpression(opts)),
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              aws.String(dynamodb.ReturnValueAllNew),
	}

	result, err := dynamoClient.UpdateItem(input)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %v", err)
	}

	return &users.User{
		Id:        result.Attributes["Id"].S,
		FirstName: result.Attributes["FirstName"].S,
		LastName:  result.Attributes["LastName"].S,
		Nickname:  result.Attributes["Nickname"].S,
		Password:  result.Attributes["Password"].S,
		Email:     result.Attributes["Email"].S,
		Country:   result.Attributes["Country"].S,
		CreatedAt: result.Attributes["CreatedAt"].S,
		UpdatedAt: result.Attributes["UpdatedAt"].S,
	}, nil
}

func generateUpdateExpression(opts *users.ModifyUserInput) string {
	var builder strings.Builder

	builder.WriteString("SET ")

	builder.WriteString("UpdatedAt = :updatedAt,")

	if opts.Country != nil {
		builder.WriteString("Country = :country,")
	}
	if opts.FirstName != nil {
		builder.WriteString("FirstName = :firstName,")
	}
	if opts.LastName != nil {
		builder.WriteString("LastName = :lastName,")
	}
	if opts.Nickname != nil {
		builder.WriteString("Nickname = :nickname,")
	}
	if opts.Email != nil {
		builder.WriteString("Email = :email,")
	}
	if opts.Password != nil {
		builder.WriteString("Password = :password,")
	}

	setCmd := builder.String()
	// Slice off the last comma
	setCmd = setCmd[:len(setCmd)-1]

	return setCmd
}

func generateExpressionAttributeValues(
	opts *users.ModifyUserInput,
) (map[string]*dynamodb.AttributeValue, error) {
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":updatedAt": {
			S: aws.String(time.Now().Format(time.RFC3339)),
		},
	}

	if opts.FirstName != nil {
		expressionAttributeValues[":firstName"] = &dynamodb.AttributeValue{
			S: opts.FirstName,
		}
	}

	if opts.LastName != nil {
		expressionAttributeValues[":lastName"] = &dynamodb.AttributeValue{
			S: opts.LastName,
		}
	}

	if opts.Nickname != nil {
		expressionAttributeValues[":nickname"] = &dynamodb.AttributeValue{
			S: opts.Nickname,
		}
	}

	if opts.Password != nil {
		hashedPassword, err := hashPassword(opts.GetPassword())
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %v", err)
		}
		expressionAttributeValues[":password"] = &dynamodb.AttributeValue{
			S: aws.String(string(hashedPassword)),
		}
	}

	if opts.Email != nil {
		expressionAttributeValues[":email"] = &dynamodb.AttributeValue{
			S: opts.Email,
		}
	}

	if opts.Country != nil {
		expressionAttributeValues[":country"] = &dynamodb.AttributeValue{
			S: opts.Country,
		}
	}

	return expressionAttributeValues, nil
}
