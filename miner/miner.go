package miner

import (
	bc "blockchainnetwork/blockchain"
	node "blockchainnetwork/node"
	"blockchainnetwork/node/proto"
	"context"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Miner struct {
	// miner identification
	name string

	// the nodes that the miner is in communication with
	nodePool node.NodeClientPool

	// data necessary to mine the next block
	dataForMining *DataForMining
}

// data necessary to mine the next block
type DataForMining struct {
	lastBlock *bc.Block
	threshold []byte
}

func MakeMiner(name string, nodePool node.NodeClientPool) *Miner {

	miner := &Miner{
		name:     name,
		nodePool: nodePool,
	}

	// TODO: start the background routine to check for new blocks in target blockchains
	client, err := miner.nodePool.GetClient("[::1]:8080")
	if err != nil {
		log.Fatal("Could not connect to client")
	}

	ctx := context.Background()
	req := &proto.GetLastBlockRequest{}

	resp, _ := client.GetLastBlock(ctx, req)

	miner.dataForMining = &DataForMining{
		lastBlock: node.ProtoBlockToBlockchainBlock(resp.LastBlock),
	}

	return miner
}

func (miner *Miner) firstCandidateBlock() *bc.Block {
	return &bc.Block{
		Index:     miner.dataForMining.lastBlock.Index + 1,
		PrevHash:  miner.dataForMining.lastBlock.PrevHash,
		Timestamp: timestamppb.Now(),
		Nonce:     uint64(0),
		Data:      "Block mined by " + miner.name,
	}
}

func (miner *Miner) nextCandidateBlock(candidateBlock *bc.Block) {
	candidateBlock.Nonce += 1
}

func (miner *Miner) Mine() {
	candidateBlock := miner.firstCandidateBlock()

	for !bc.HashSatisfiesThreshold(candidateBlock, miner.dataForMining.threshold) {
		miner.nextCandidateBlock(candidateBlock)
	}

	// send it to the blockchain
	miner.sendBlockToNode("[::1]:8080", candidateBlock)
}

func (miner *Miner) sendBlockToNode(nodeName string, block *bc.Block) {
	client, err := miner.nodePool.GetClient(nodeName)
	if err != nil {
		log.Fatal("Could not connect to client")
	}

	// choose the parameters
	ctx := context.Background()
	req := &proto.AppendBlocksRequest{
		Blocks: node.BlockchainBlocksToProtoBlocks([]*bc.Block{block}),
	}

	resp, err := client.AppendBlocks(ctx, req)
	if err != nil {
		log.Println("Error contacting node in sendBlockToNode")
		return
	}
	if !resp.Success {
		log.Println("Block rejected by node")
	}

	//fmt.Printf("resp.Success: %d, resp.LastBlockHash: %b, resp.LastBlockIndex\n", resp.Success, resp.LastBlockHash, resp.LastBlockIndex)
}
