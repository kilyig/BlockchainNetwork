package test

import (
	"context"
	"sync"

	node "blockchainnetwork/node"
	"blockchainnetwork/node/proto"

	"google.golang.org/grpc"
)

type MockNodeClientPool struct {
	mu    sync.RWMutex
	nodes map[string]*MockNodeClient
}

type MockNodeClient struct {
	inMemNode *node.Node
}

func (cp *MockNodeClientPool) GetClient(nodeName string) (proto.NodeClient, error) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	return cp.nodes[nodeName].inMemNode, nil
}

func MakeMockNodeClientPool(nodes []string) *MockNodeClientPool {
	cp := &MockNodeClientPool{
		nodes: make(map[string]*MockNodeClient),
	}

	for _, nodeName := range nodes {
		client, err := makeMockNodeClient(nodeName, cp)
		if err == nil {
			cp.nodes[nodeName] = client
		}
	}

	return cp
}

func makeMockNodeClient(addr string, nodePool node.NodeClientPool) (*MockNodeClient, error) {
	return &MockNodeClient{
		inMemNode: node.MakeNode(addr, nodePool),
	}, nil
}

func (nc *MockNodeClient) GetBlocks(ctx context.Context, req *proto.GetBlocksRequest, opts ...grpc.CallOption) (*proto.GetBlocksResponse, error) {
	return nc.inMemNode.GetBlocks(ctx, req, opts...)
}

func (nc *MockNodeClient) AppendBlocks(ctx context.Context, req *proto.AppendBlocksRequest, opts ...grpc.CallOption) (*proto.AppendBlocksResponse, error) {
	return nc.inMemNode.AppendBlocks(ctx, req, opts...)
}

func (nc *MockNodeClient) GetLastBlock(ctx context.Context, req *proto.GetLastBlockRequest, opts ...grpc.CallOption) (*proto.GetLastBlockResponse, error) {
	return nc.inMemNode.GetLastBlock(ctx, req, opts...)
}
