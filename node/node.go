package node

import (
	"blockchainnetwork/blockchain"
	bc "blockchainnetwork/blockchain"
	proto "blockchainnetwork/node/proto"
	"bytes"
	"context"
	"log"
	"time"
)

const (
	daemonTimeDelta = 5 * time.Second // for the ticker
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

	// start the background routine to check for new blocks in other nodes
	go node.syncBlockchainDaemon()

	return node
}

func (n *Node) addNode(nodeName string) {
	n.nodes[nodeName] = struct{}{}
}

func (n *Node) syncBlockchainDaemon() {
	for {
		time.Sleep(daemonTimeDelta)
		n.syncWithNetwork()
	}
}

func (n *Node) syncWithNetwork() {
	for node := range n.nodes {
		go func(nodeName string) {
			n.syncWithNode(nodeName)
		}(node)
	}
}

func (n *Node) syncWithNode(nodeName string) {
	client, err := n.nodePool.GetClient(nodeName)
	if err != nil {
		log.Fatal("Could not connect to client")
	}

	// prepare the request
	ctx := context.Background()
	allBlocks, err := n.blockchain.GetBlocks(1)
	if err != nil {
		return
	}
	req := &proto.AddBlocksRequest{
		Blocks:         BlockchainBlocksToProtoBlocks(allBlocks),
		PrevBlockIndex: 0,
		PrevBlockHash:  bc.HashBlock(n.blockchain.GetBlock(0)),
	}

	// send the RPC and process the reply
	resp, err := client.AddBlocks(ctx, req)
	if err == nil {
		n.handleAddBlocksResponse(resp)
	}
}

func (n *Node) handleAddBlocksResponse(resp *proto.AddBlocksResponse) {

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

func (n *Node) AddBlocks(ctx context.Context, req *proto.AddBlocksRequest) (*proto.AddBlocksResponse, error) {
	lastAddedBlockIndex, ok := n.blockchain.AddBlocks(ProtoBlocksToBlockchainBlocks(req.Blocks))

	if ok && len(req.Blocks) != 0 {
		log.Printf("Added blocks from #%d to #%d with data: %s\n", req.GetBlocks()[0].Index, lastAddedBlockIndex, req.GetBlocks()[0].Data)
	}

	// the request is successfull if this blockchain agrees with the Prev fields
	// in the request
	success := false
	agreementBlock := n.blockchain.GetBlock(req.PrevBlockIndex)
	if agreementBlock != nil {
		success = bytes.Equal(req.PrevBlockHash, bc.HashBlock(agreementBlock))
	}

	newLastBlock := n.blockchain.LastBlock()
	return &proto.AddBlocksResponse{
		LastBlockIndex: newLastBlock.Index,
		LastBlockHash:  bc.HashBlock(newLastBlock),
		Success:        success,
	}, nil
}

func (n *Node) GetLastBlock(ctx context.Context, req *proto.GetLastBlockRequest) (*proto.GetLastBlockResponse, error) {
	return &proto.GetLastBlockResponse{
		LastBlock: BlockchainBlockToProtoBlock(n.blockchain.LastBlock()),
	}, nil
}
