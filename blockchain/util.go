package blockchain

import (
	"crypto/sha256"
	"encoding/json"
)

func SHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func hashBlock(block *Block) ([]byte, error) {
	blockBytes, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	return SHA256(blockBytes), nil
}

func isValidBlock(block *Block) bool {
	return true
}

// https://gist.github.com/miguelmota/3dee93d8b7340e33fc474eb3abb7d450
