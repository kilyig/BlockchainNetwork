package miner

import (
	bc "blockchainnetwork/blockchain"
	node "blockchainnetwork/node"
	"blockchainnetwork/node/proto"
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	daemonTimeDelta = 5 * time.Second // for the ticker
)

type Miner struct {
	mu sync.RWMutex

	mineDelay time.Duration

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

func MakeMiner(name string, nodePool node.NodeClientPool, nodes []string, mineDelay uint64) *Miner {

	miner := &Miner{
		name:          name,
		nodePool:      nodePool,
		nodes:         make(map[string]struct{}),
		mineDelay:     time.Duration(mineDelay) * time.Nanosecond,
		dataForMining: nil,
	}

	// add the nodes to the local registry
	for _, neighborNode := range nodes {
		miner.addNode(neighborNode)
	}

	miner.collectDataForMiningFromNetwork()
	go miner.collectDataForMiningFromNetworkDaemon()

	return miner
}

func (miner *Miner) addNode(nodeName string) {
	miner.nodes[nodeName] = struct{}{}
}

func (miner *Miner) firstCandidateBlock() (*bc.Block, error) {
	if miner.dataForMining == nil {
		return nil, status.Error(codes.NotFound, "miner does not have data to mine")
	}

	return &bc.Block{
		Index:     miner.dataForMining.lastBlockIndex + 1,
		PrevHash:  miner.dataForMining.lastBlockHash,
		Timestamp: timestamppb.Now(),
		Nonce:     uint64(0),
		Data:      "Block mined by " + miner.name,
	}, nil
}

func (miner *Miner) collectDataForMiningFromNetworkDaemon() {
	for {
		time.Sleep(daemonTimeDelta)
		miner.collectDataForMiningFromNetwork()
	}
}

func (miner *Miner) collectDataForMiningFromNetwork() {
	for node := range miner.nodes {
		go func(nodeName string) {
			miner.getDataForMiningFromNode(nodeName)
		}(node)
	}
}

func (miner *Miner) getDataForMiningFromNode(nodeName string) {
	client, err := miner.nodePool.GetClient(nodeName)
	if err != nil {
		log.Fatal("Could not connect to client")
	}

	ctx := context.Background()
	req := &proto.AddBlocksRequest{
		Blocks: make([]*proto.Block, 0),
	}

	resp, err := client.AddBlocks(ctx, req)
	if err == nil {
		miner.handleAddBlocksResponse(resp)
	}
}

func (miner *Miner) nextCandidateBlock(candidateBlock *bc.Block) {
	candidateBlock.Nonce += 1
}

func (miner *Miner) MineContinuously() {
	// i := 5
	// k := 0
	for /*k < i*/ {
		miner.MineOneBlock()
		// k += 1
	}
}

func (miner *Miner) MineOneBlock() {
	candidateBlock, err := miner.firstCandidateBlock()
	if err != nil {
		return
	}

	log.Printf("(miner: %s) Mining a block with index %d", miner.name, candidateBlock.Index)

	for !bc.BlockHashSatisfiesThreshold(candidateBlock, miner.dataForMining.threshold) {
		time.Sleep(miner.mineDelay)
		miner.nextCandidateBlock(candidateBlock)
	}
	//fmt.Printf("valid hash: %x\n", bc.HashBlock(candidateBlock))

	// send it to every node in the network
	miner.sendBlockToNetwork(candidateBlock)
}

func (miner *Miner) sendBlockToNetwork(block *bc.Block) {
	wg := new(sync.WaitGroup)

	for node := range miner.nodes {
		wg.Add(1)
		go func(nodeName string) {
			miner.sendBlockToNode(nodeName, block)
			wg.Done()
		}(node)
	}

	wg.Wait()
}

func (miner *Miner) sendBlockToNode(nodeName string, block *bc.Block) {
	client, err := miner.nodePool.GetClient(nodeName)
	if err != nil {
		log.Fatal("Could not connect to client")
	}

	// choose the parameters
	ctx := context.Background()
	req := &proto.AddBlocksRequest{
		Blocks: node.BlockchainBlocksToProtoBlocks([]*bc.Block{block}),
	}

	resp, err := client.AddBlocks(ctx, req)
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
	miner.handleAddBlocksResponse(resp)
}

func (miner *Miner) handleAddBlocksResponse(resp *proto.AddBlocksResponse) {
	miner.mu.Lock()
	defer miner.mu.Unlock()

	// if the blockchain is longer than what you think, continue building on
	// that blockchain
	if miner.dataForMining == nil || resp.LastBlockIndex > miner.dataForMining.lastBlockIndex {
		miner.dataForMining = &DataForMining{
			lastBlockIndex: resp.LastBlockIndex,
			lastBlockHash:  resp.LastBlockHash,
			threshold:      []byte{0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		}
	}
	// TODO: will be "miner.dataForMining.threshold = resp.Threshold" after updating the .proto file
}
