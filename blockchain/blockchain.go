package blockchain

import (
	"bytes"
	"time"
)

type Blockchain struct {
	Blocks []*Block
}

type Block struct {
	Index     uint64
	PrevHash  []byte
	Timestamp time.Time
	Data      string
}

func MakeBlockchain() *Blockchain {
	// the initial blockchain has only the genesis block
	blocks := make([]*Block, 0)
	blocks = append(blocks, &Block{
		Index:     0,
		PrevHash:  make([]byte, 32), // SHA256 outputs have 32 bytes
		Timestamp: time.Now(),
		Data:      "",
	})
	return &Blockchain{
		Blocks: blocks,
	}
}

func (bc *Blockchain) lastBlock() *Block {
	return bc.Blocks[len(bc.Blocks)]
}

/* For a block to be added to the blockchain,
 * 1) its index should be one more than the last block's index
 * 2) its "prevHash" field should be equal to the last block's hash
 * 3) its hash should be valid
 */
func (bc *Blockchain) isValidNextBlock(candidateBlock *Block) bool {
	lastBlock := bc.lastBlock()

	if candidateBlock.Index == lastBlock.Index+1 &&
		bytes.Equal(candidateBlock.PrevHash, hashBlock(lastBlock)) &&
		bc.hasValidHash(candidateBlock) {
		return true
	}
	return false
}

func (bc *Blockchain) hasValidHash(block *Block) bool {
	return true
}

func (bc *Blockchain) addBlock(newBlock *Block) (uint64, bool) {
	if bc.isValidNextBlock(newBlock) {
		bc.Blocks = append(bc.Blocks, newBlock)
		return uint64(len(bc.Blocks)), true
	}
	return uint64(len(bc.Blocks)), false
}
