package miner

import (
	bc "blockchainnetwork/blockchain"
	"time"
)

type Miner struct {
	// miner identificaiton
	name string

	// the blockchain that the miner mines blocks for
	blockchain *bc.Blockchain

	// update when the target blockchain receives a new block
	currentIndex    uint64
	currentPrevHash []byte
	threshold       []byte
}

func (miner *Miner) mine() {
	minedBlock := &bc.Block{
		Index:     miner.currentIndex,
		PrevHash:  miner.currentPrevHash,
		Timestamp: time.Now(),
		Nonce:     uint64(0),
		Data:      "Block mined by " + miner.name,
	}

	for !bc.hashSatisfiesThreshold(minedBlock, miner.threshold) {
		minedBlock.Nonce += 1
	}

	// send it to the blockchain

}
