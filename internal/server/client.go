package server

import (
	"fmt"
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/litetable/litetable-db/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	conn   *grpc.ClientConn
	client proto.LitetableServiceClient

	rpcConnString string
}

// NewClient creates a new LiteTable gRPC client
func NewClient() (*GrpcClient, error) {
	// TODO: allow the setting of secure credentials
	serverAddress, err := litetable.GetFromConfig(litetable.ServerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get server address: %w", err)
	}

	serverRPCPort, err := litetable.GetFromConfig(litetable.ServerRPCPort)
	if err != nil {
		return nil, fmt.Errorf("failed to get server RPC port: %w", err)
	}

	connString := fmt.Sprintf("%s:%s", serverAddress, serverRPCPort)
	conn, err := grpc.NewClient(connString,
		grpc.WithTransportCredentials(insecure.
			NewCredentials()))
	if err != nil {
		return nil, err
	}

	ltClient := proto.NewLitetableServiceClient(conn)
	return &GrpcClient{
		rpcConnString: connString,
		conn:          conn,
		client:        ltClient,
	}, nil
}

// Close closes the client connection
func (g *GrpcClient) Close() error {
	if err := g.conn.Close(); err != nil {
		return err
	}

	return nil
}
