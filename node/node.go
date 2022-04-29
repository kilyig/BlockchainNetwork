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

	logrus.Info("A group of walrus emerges from the ocean")

	added_all_blocks := true
	for _, newBlock := range req.GetBlocks() {
		if n.isValidNextBlock(newBlock) {
			n.blockchain.AddBlock(ProtoBlockToBlockchainBlock(newBlock))
		} else {
			added_all_blocks = false
			break
		}
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
