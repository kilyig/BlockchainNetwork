package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
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

// checks if the hash of block is less than or equal to the threshold
func HashSatisfiesThreshold(block *Block, threshold []byte) bool {
	return bytes.Compare(HashBlock(block), threshold) != 1
}

// https://gist.github.com/miguelmota/3dee93d8b7340e33fc474eb3abb7d450
