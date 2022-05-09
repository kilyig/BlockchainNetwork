package node

import (
	"blockchainnetwork/blockchain"
	bc "blockchainnetwork/blockchain"
	proto "blockchainnetwork/node/proto"
	"context"

	"github.com/sirupsen/logrus"
)

type Node struct {
	proto.UnimplementedNodeServer

	// node identification
	name string

	// connections to nodes that this node is in communication with
	nodePool NodeClientPool
	nodes    map[string]struct{}

	// this node's blockchain
	blockchain *blockchain.Blockchain
}

func MakeNode(name string, nodePool NodeClientPool, nodes []string) *Node {

	node := &Node{
		name:       name,
		nodePool:   nodePool,
		blockchain: blockchain.MakeBlockchain(),
		nodes:      make(map[string]struct{}),
	}

	// add the nodes to the local registry
	for _, neighborNode := range nodes {
		node.addNode(neighborNode)
	}

	// TODO: start the background routine to check for new blocks in other nodes

	return node
}

func (n *Node) addNode(nodeName string) {
	n.nodes[nodeName] = struct{}{}
}

func (n *Node) GetBlocks(ctx context.Context, req *proto.GetBlocksRequest) (*proto.GetBlocksResponse, error) {
	blocks, err := n.blockchain.GetBlocks(req.GetFirstBlockIndex())

	if err != nil {
		return &proto.GetBlocksResponse{
			Blocks: nil,
		}, err
	}

	return &proto.GetBlocksResponse{
		Blocks: BlockchainBlocksToProtoBlocks(blocks),
	}, nil
}

func (n *Node) isValidNextBlock(block *proto.Block) bool {
	return n.blockchain.IsValidNextBlock(ProtoBlockToBlockchainBlock(block))
}

func (n *Node) AppendBlocks(ctx context.Context, req *proto.AppendBlocksRequest) (*proto.AppendBlocksResponse, error) {

	// logrus.WithFields(logrus.Fields{
	// 	"animal": "walrus",
	// 	"size":   10,
	// }).Info("A group of walrus emerges from the ocean")

	added_all_blocks := true
	for _, newBlock := range req.GetBlocks() {
		if n.isValidNextBlock(newBlock) {
			n.blockchain.AddBlock(ProtoBlockToBlockchainBlock(newBlock))
		} else {
			added_all_blocks = false
			break
		}
	}

	if added_all_blocks && len(req.GetBlocks()) != 0 {
		logrus.Infof("Added block #%d with data: %s\n", req.GetBlocks()[0].Index, req.GetBlocks()[0].Data)
	}

	newLastBlock := n.blockchain.LastBlock()
	return &proto.AppendBlocksResponse{
		LastBlockIndex: newLastBlock.Index,
		LastBlockHash:  bc.HashBlock(newLastBlock),
		Success:        added_all_blocks,
	}, nil

}

func (n *Node) GetLastBlock(ctx context.Context, req *proto.GetLastBlockRequest) (*proto.GetLastBlockResponse, error) {
	return &proto.GetLastBlockResponse{
		LastBlock: BlockchainBlockToProtoBlock(n.blockchain.LastBlock()),
	}, nil
}
