package test

import (
	"sync"

	node "blockchainnetwork/node"
	"blockchainnetwork/node/proto"
)

type MockNodeClientPool struct {
	mu    sync.RWMutex
	nodes map[string]*MockNodeClient
}

type MockNodeClient struct {
	inMemNode *node.Node
}

func MakeMockNodeClientPool(nodes []string) *MockNodeClientPool {
	cp := &MockNodeClientPool{
		nodes: make(map[string]*MockNodeClient),
	}

	for _, node := range nodes {
		cp.nodes[node]
	}
}

func makeMockNodeClient(addr string) (*MockNodeClient, error) {
	return &MockNodeClient{
		inMemNode: node.MakeNode(addr),
	}
}

func (cp *MockNodeClientPool) GetClient(nodeName string) (proto.NodeClient, error) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	return cp.nodes[nodeName].inMemNode, nil
}
