package miner

import (
	bc "blockchainnetwork/blockchain"
	node "blockchainnetwork/node"
	"blockchainnetwork/node/proto"
	"context"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	daemonTimeDelta = 5 * time.Second // for the ticker
)

type Miner struct {
	// miner identification
	name string

	// the nodes that the miner is in communication with
	nodePool node.NodeClientPool
	nodes    map[string]struct{}

	// data necessary to mine the next block
	dataForMining *DataForMining
}

// data necessary to mine the next block
type DataForMining struct {
	lastBlockIndex uint64
	lastBlockHash  []byte
	threshold      []byte
}

func MakeMiner(name string, nodePool node.NodeClientPool, nodes []string) *Miner {

	miner := &Miner{
		name:     name,
		nodePool: nodePool,
		nodes:    make(map[string]struct{}),
	}

	// add the nodes to the local registry
	for _, neighborNode := range nodes {
		miner.addNode(neighborNode)
	}

	miner.updateDataForMining()
	go miner.updateDataForMiningDaemon()

	return miner
}

func (miner *Miner) addNode(nodeName string) {
	miner.nodes[nodeName] = struct{}{}
}

func (miner *Miner) firstCandidateBlock() *bc.Block {
	return &bc.Block{
		Index:     miner.dataForMining.lastBlockIndex + 1,
		PrevHash:  miner.dataForMining.lastBlockHash,
		Timestamp: timestamppb.Now(),
		Nonce:     uint64(0),
		Data:      "Block mined by " + miner.name,
	}
}

func (miner *Miner) updateDataForMiningDaemon() {
	for {
		time.Sleep(daemonTimeDelta)
		miner.updateDataForMining()
	}
}

func (miner *Miner) updateDataForMining() {
	client, err := miner.nodePool.GetClient("[::1]:8080")
	if err != nil {
		log.Fatal("Could not connect to client")
	}

	ctx := context.Background()
	req := &proto.AppendBlocksRequest{
		Blocks: make([]*proto.Block, 0),
	}
	resp, err := client.AppendBlocks(ctx, req)

	if err == nil {
		miner.handleAppendBlocksResponse(resp)
	}
}

func (miner *Miner) nextCandidateBlock(candidateBlock *bc.Block) {
	candidateBlock.Nonce += 1
}

func (miner *Miner) MineContinuously() {
	for {
		miner.MineOneBlock()
	}
}

func (miner *Miner) MineOneBlock() {
	candidateBlock := miner.firstCandidateBlock()

	log.Printf("(miner: %s) Mining a block with index %d", miner.name, candidateBlock.Index)

	for !bc.BlockHashSatisfiesThreshold(candidateBlock, miner.dataForMining.threshold) {
		miner.nextCandidateBlock(candidateBlock)
	}
	//fmt.Printf("valid hash: %x\n", bc.HashBlock(candidateBlock))

	// send it to every node in the network
	miner.sendBlockToNetwork(candidateBlock)
}

func (miner *Miner) sendBlockToNetwork(block *bc.Block) {
	for node := range miner.nodes {
		miner.sendBlockToNode(node, block)
	}
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
	if resp.Success {
		log.Printf("(miner: %s) Block accepted by node %s\n", miner.name, nodeName)
	} else {
		log.Printf("(miner: %s) Block rejected by node %s\n", miner.name, nodeName)
	}

	//update your mining data no matter what
	miner.handleAppendBlocksResponse(resp)
}

func (miner *Miner) handleAppendBlocksResponse(resp *proto.AppendBlocksResponse) {
	miner.dataForMining = &DataForMining{
		lastBlockIndex: resp.LastBlockIndex,
		lastBlockHash:  resp.LastBlockHash,
		threshold:      []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	// TODO: will be "miner.dataForMining.threshold = resp.Threshold" after updating the .proto file
}
