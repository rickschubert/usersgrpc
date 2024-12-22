package server

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rickschubert/usersgrpc/config"
	"github.com/rickschubert/usersgrpc/server/internal/handlers"
	"github.com/rickschubert/usersgrpc/users"
	"go.uber.org/zap"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type server struct {
	logger       *zap.SugaredLogger
	dynamoClient *dynamodb.DynamoDB
	users.UnimplementedUsersServer
}

func New(
	logger *zap.SugaredLogger,
	dynamoClient *dynamodb.DynamoDB,
	healthServer *health.Server,
) users.UsersServer {
	healthServer.SetServingStatus(config.ServiceName(), grpc_health_v1.HealthCheckResponse_SERVING)

	return &server{
		dynamoClient: dynamoClient,
		logger:       logger,
	}
}

func (s server) ListUsers(
	ctx context.Context,
	opts *users.ListUsersInput,
) (*users.ListUsersResponse, error) {
	defer recoverFromPanics()
	return handlers.ListUsers(ctx, s.logger, s.dynamoClient, opts)
}

func (s server) AddUser(ctx context.Context, opts *users.AddUserInput) (*users.User, error) {
	defer recoverFromPanics()
	return handlers.AddUser(ctx, s.logger, s.dynamoClient, opts)
}

func (s server) ModifyUser(ctx context.Context, opts *users.ModifyUserInput) (*users.User, error) {
	defer recoverFromPanics()
	return handlers.ModifyUser(ctx, s.logger, s.dynamoClient, opts)
}

func (s server) RemoveUser(ctx context.Context, opts *users.RemoveUserInput) (*users.User, error) {
	defer recoverFromPanics()
	return handlers.RemoveUser(ctx, s.logger, s.dynamoClient, opts)
}

func (s server) GetUser(ctx context.Context, opts *users.GetUserInput) (*users.User, error) {
	defer recoverFromPanics()
	return handlers.GetUser(ctx, s.logger, s.dynamoClient, opts)
}

func recoverFromPanics() {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic: %v", r)
	}
}
