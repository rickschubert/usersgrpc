package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rickschubert/usersgrpc/config"
)

func New() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(config.Region()),
		Endpoint: aws.String(config.DynamoDBEndpoint()),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return dynamodb.New(sess), nil
}
