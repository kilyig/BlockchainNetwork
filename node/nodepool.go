package node

import (
	"sync"

	"blockchainnetwork/node/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NodeClientPool interface {
	GetClient(nodeName string) (proto.NodeClient, error)
}

type GRPCNodeClientPool struct {
	mutex   sync.RWMutex
	clients map[string]proto.NodeClient
}

func MakeGRPCClientPool(nodes []string) *GRPCNodeClientPool {

	gRPCClientPool := &GRPCNodeClientPool{
		clients: make(map[string]proto.NodeClient),
	}

	// set up connections with the nodes with Dial()
	for _, node := range nodes {
		client, err := makeNodeClient(node)
		if err != nil {
			return nil
		}
		gRPCClientPool.clients[node] = client
	}

	return gRPCClientPool
}

func makeNodeClient(addr string) (proto.NodeClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	channel, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	return proto.NewNodeClient(channel), nil
}
