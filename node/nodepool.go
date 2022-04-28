package node

import (
	"context"
	"sync"

	"blockchainnetwork/node/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type NodeClientPool interface {
	GetClient(nodeName string) (proto.NodeClient, error)
}

type GRPCNodeClientPool struct {
	mu      sync.RWMutex
	clients map[string]proto.NodeClient
}

func MakeGRPCNodeClientPool(nodes []string) *GRPCNodeClientPool {

	gRPCClientPool := &GRPCNodeClientPool{
		clients: make(map[string]proto.NodeClient),
	}

	// set up connections with the nodes with Dial()
	for _, node := range nodes {
		client, err := makeGRPCNodeClient(node)
		if err != nil {
			return nil
		}
		gRPCClientPool.clients[node] = client
	}

	return gRPCClientPool
}

func makeGRPCNodeClient(addr string) (proto.NodeClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	channel, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	return proto.NewNodeClient(channel), nil
}

func (cp *GRPCNodeClientPool) GetClient(nodeName string) (proto.NodeClient, error) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	client, ok := cp.clients[nodeName]
	if !ok {
		return nil, status.Error(codes.NotFound, "error while getting the client.")
	}
	return client, nil
}

func (cp *GRPCNodeClient) GetBlocks(ctx context.Context, in *GetBlocksRequest, opts ...grpc.CallOption) (*GetBlocksResponse, error)
