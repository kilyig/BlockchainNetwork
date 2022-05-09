package blockchain

import (
	"bytes"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type Blockchain struct {
	Blocks    []*Block
	threshold []byte
}

type Block struct {
	Index     uint64
	PrevHash  []byte
	Timestamp *timestamppb.Timestamp // TODO: will be changed to time.Time in the future
	Nonce     uint64
	Data      string
}

func MakeBlockchain() *Blockchain {
	// the initial blockchain has only the genesis block
	blocks := make([]*Block, 0)
	blocks = append(blocks, &Block{
		Index:     0,
		PrevHash:  make([]byte, 32), // SHA256 outputs have 32 bytes
		Timestamp: timestamppb.Now(),
		Data:      "",
	})
	return &Blockchain{
		Blocks:    blocks,
		threshold: []byte{0, 0, 127, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}

func (bc *Blockchain) GetBlocks(firstBlockIndex uint64) ([]*Block, error) {
	if bc.LastBlock().Index < firstBlockIndex {
		return nil, status.Error(codes.OutOfRange, "this blockchain does not have a block with the requested index.")
	}
	return bc.Blocks[firstBlockIndex:], nil
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

/*
 * For a block to be added to the blockchain,
 * 1) its index should be one more than the last block's index
 * 2) its "prevHash" field should be equal to the last block's hash
 * 3) its hash should be valid
 */
func (bc *Blockchain) IsValidNextBlock(candidateBlock *Block) bool {
	lastBlock := bc.LastBlock()

	if candidateBlock.Index == lastBlock.Index+1 &&
		bytes.Equal(candidateBlock.PrevHash, HashBlock(lastBlock)) &&
		bc.HasValidHash(candidateBlock) {
		return true
	}
	return false
}

func (bc *Blockchain) HasValidHash(block *Block) bool {
	return BlockHashSatisfiesThreshold(block, bc.threshold)
}

func (bc *Blockchain) AddBlock(newBlock *Block) (uint64, bool) {
	if bc.IsValidNextBlock(newBlock) {
		bc.Blocks = append(bc.Blocks, newBlock)
		return uint64(len(bc.Blocks)), true
	}
	return uint64(len(bc.Blocks)), false
}
