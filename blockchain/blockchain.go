package blockchain

import (
	"time"
)

type Blockchain struct {
	Size   uint64
	Blocks []*Block
}

type Block struct {
	Version       uint64
	HashPrevBlock []byte
	Timestamp     time.Time
}
