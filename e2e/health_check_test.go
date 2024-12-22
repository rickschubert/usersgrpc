package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rickschubert/usersgrpc/config"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestHealthCheck(t *testing.T) {
	healthConn, err := grpc.NewClient(
		fmt.Sprintf(":%d", config.Port()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal("fail to create grpc client: %w", err)
	}
	defer healthConn.Close()

	healthClient := grpc_health_v1.NewHealthClient(healthConn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	req := &grpc_health_v1.HealthCheckRequest{Service: config.ServiceName()}
	resp, err := healthClient.Check(ctx, req)

	require.NoError(t, err)
	require.Equal(t, "SERVING", resp.Status.String())
}
