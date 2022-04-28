package miner

import (
	bc "blockchainnetwork/blockchain"
	node "blockchainnetwork/node"

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

func (miner *Miner) mine() {
	candidateBlock := miner.firstCandidateBlock()

	for !bc.HashSatisfiesThreshold(candidateBlock, miner.dataForMining.threshold) {
		miner.nextCandidateBlock(candidateBlock)
	}

	// send it to the blockchain

}
