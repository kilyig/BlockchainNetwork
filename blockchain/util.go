package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

func SHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func HashBlock(block *Block) []byte {
	blockBytes, err := json.Marshal(block)
	if err != nil {
		return nil
	}
	return SHA256(blockBytes)
}

func HashSatisfiesThreshold(data []byte, threshold []byte) bool {
	return bytes.Compare(data, threshold) < 1
}

// checks if the hash of block is less than the threshold
func BlockHashSatisfiesThreshold(block *Block, threshold []byte) bool {
	return bytes.Compare(HashBlock(block), threshold) != 1
}

func PrintBlockchainBlock(block *Block) {
	fmt.Printf("Index: %d\n", block.Index)
	fmt.Printf("	PrevHash: %d \n", block.PrevHash)
	fmt.Printf("	Timestamp: %d\n", block.Timestamp.Seconds)
	fmt.Printf("	Nonce: %d\n", block.Nonce)
	fmt.Printf("	Data: %s\n", block.Data)
}

// https://gist.github.com/miguelmota/3dee93d8b7340e33fc474eb3abb7d450
