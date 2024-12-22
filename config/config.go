package config

import (
	"fmt"
	"os"

	"github.com/samber/lo"
)

func ServiceName() string {
	return "usersgrpc"
}

func TableName() string {
	return "users"
}

func Port() int {
	return 50051
}

func Region() string {
	return "us-east-1"
}

func DynamoDBEndpoint() string {
	return fmt.Sprintf(
		"http://%s:4566",
		lo.Ternary(os.Getenv("LOCAL") == "true", "localhost", "localstack"),
	)
}

func PaginationSize() int {
	return 2
}
