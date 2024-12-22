package main

import (
	"fmt"
	"log"
	"net"

	"github.com/rickschubert/usersgrpc/config"
	"github.com/rickschubert/usersgrpc/db"
	"github.com/rickschubert/usersgrpc/server"
	"github.com/rickschubert/usersgrpc/users"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	logger := setupLogger()
	defer logger.Sync()

	grpcServer := grpc.NewServer()
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	dynamoClient, err := db.New()
	if err != nil {
		logger.Fatalf("failed to create dynamoDB client: %v", err)
	}

	usersServer := server.New(logger, dynamoClient, healthServer)

	users.RegisterUsersServer(grpcServer, usersServer)

	logger.Info("Starting server on port ", config.Port())
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port()))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	grpcServer.Serve(lis)
}

func setupLogger() *zap.SugaredLogger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		log.Fatalf("failed to create zap logger from config: %v", err)
	}
	return logger.Sugar()
}
