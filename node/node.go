package node

import (
	"blockchainnetwork/blockchain"
	"blockchainnetwork/node/proto"
	"context"

	"google.golang.org/grpc"
)

type Node struct {
	// node identification
	name string

	// the nodes that this node is in communication with
	nodePool NodeClientPool

	// this node's blockchain
	blockchain *blockchain.Blockchain
}

func MakeNode(name string, nodePool NodeClientPool) *Node {

	node := &Node{
		name:       name,
		nodePool:   nodePool,
		blockchain: blockchain.MakeBlockchain(),
	}

	// TODO: start the background routine to check for new blocks in target blockchains

	return node
}

func (n *Node) GetBlocks(ctx context.Context, in *proto.GetBlocksRequest, opts ...grpc.CallOption) (*proto.GetBlocksResponse, error) {
	panic("AAAA")
}

func (n *Node) AppendBlocks(ctx context.Context, in *proto.AppendBlocksRequest, opts ...grpc.CallOption) (*proto.AppendBlocksResponse, error) {
	panic("AAAA")
}

func (n *Node) GetLastBlock(ctx context.Context, in *proto.GetLastBlockRequest, opts ...grpc.CallOption) (*proto.GetLastBlockResponse, error) {
	panic("AAAA")
}
