package blockchain

import (
	"time"
)

type Blockchain struct {
	Size   uint64
	Blocks []*Block
}

type Block struct {
	Version   uint64
	Index     uint64
	PrevHash  []byte
	Timestamp time.Time
	Data      string
}

func MakeBlockchain() *Blockchain {
	return &Blockchain{
		Size:   0,
		Blocks: make([]*Block, 0),
	}
}

func (bc *Blockchain) addBlock(block *Block) (uint64, error) {
	panic("unimplemented")
}
