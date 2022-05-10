package blockchain

import (
	"bytes"
	"time"

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

	blockchain := &Blockchain{
		Blocks:    make([]*Block, 0),
		threshold: []byte{0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	// the initial blockchain has only the genesis block
	genesisBlock := &Block{
		Index:     0,
		PrevHash:  make([]byte, 32), // SHA256 outputs have 32 bytes
		Timestamp: timestamppb.New(time.Time{}),
		Data:      "",
		Nonce:     0,
	}
	blockchain.findValidNonce(genesisBlock, blockchain.threshold)
	blockchain.Blocks = append(blockchain.Blocks, genesisBlock)

	return blockchain
}

func (bc *Blockchain) findValidNonce(block *Block, threshold []byte) {
	for !BlockHashSatisfiesThreshold(block, threshold) {
		block.Nonce += 1
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

func (bc *Blockchain) IsValidBlock(candidateBlock *Block) bool {
	// the block is invalid if its index is higher than the index of the next
	// block that the blockchain can currently accept
	if candidateBlock.Index > bc.LastBlock().Index+1 {
		return false
	}

	prevBlock := bc.GetBlock(candidateBlock.Index - 1)

	if candidateBlock.Index == prevBlock.Index+1 &&
		bytes.Equal(candidateBlock.PrevHash, HashBlock(prevBlock)) &&
		bc.HasValidHash(candidateBlock) {
		return true
	}
	return false
}

func (bc *Blockchain) GetBlock(index uint64) *Block {
	return bc.Blocks[index]
}

func (bc *Blockchain) HasValidHash(block *Block) bool {
	return BlockHashSatisfiesThreshold(block, bc.threshold)
}

func (bc *Blockchain) AppendBlock(newBlock *Block) (uint64, bool) {
	if bc.IsValidNextBlock(newBlock) {
		bc.Blocks = append(bc.Blocks, newBlock)
		return uint64(len(bc.Blocks)), true
	}
	return uint64(len(bc.Blocks)), false
}

func (bc *Blockchain) AddBlocks(newBlocks []*Block) (uint64, bool) {
	if len(newBlocks) == 0 {
		return bc.LastBlock().Index, true
	}

	// find the last block that can be added
	validBlockFound := false
	lastValidBlockIndex := newBlocks[0].Index
	for _, block := range newBlocks {
		if bc.IsValidBlock(block) {
			validBlockFound = true
			lastValidBlockIndex = block.Index
		}
	}

	// proceed only if (a) there was at least one valid block
	// (b) the valid blocks, when added to the blockchain, make
	// the blockchain longer than what it currently is
	if !validBlockFound || lastValidBlockIndex <= bc.LastBlock().Index {
		return bc.LastBlock().Index, false
	}

	// trim the blockchain and add the new ones
	firstNewBlockIndex := newBlocks[0].Index
	bc.Blocks = bc.Blocks[:firstNewBlockIndex]
	bc.Blocks = append(bc.Blocks, newBlocks...)

	return bc.LastBlock().Index, true
}
