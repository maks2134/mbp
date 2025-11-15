package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConfig struct {
	Address string
	Timeout time.Duration
}

func NewConnection(ctx context.Context, config ClientConfig) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", config.Address, err)
	}

	return conn, nil
}
