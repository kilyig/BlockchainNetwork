package test

import (
	"encoding/json"
	"fmt"
	"testing"

	bc "blockchainnetwork/blockchain"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPrintBlockchainBlock(t *testing.T) {
	block := &bc.Block{
		Index:     25,
		PrevHash:  make([]byte, 32),
		Timestamp: timestamppb.Now(),
		Nonce:     4566,
		Data:      "Bunu yazan tosun",
	}

	bc.PrintBlockchainBlock(block)
}

func TestHash(t *testing.T) {
	dummyWord := "dummy"
	hash := bc.SHA256([]byte(dummyWord))
	fmt.Printf("%x\n", hash)
}

func TestHashSatisfiesThreshold(t *testing.T) {
	lowestThreshold := make([]byte, 32)

	// check that nothing satisfies if threshold = 0
	dummyWord := "dummy"
	dummyWordBytes, err := json.Marshal(dummyWord)
	if err == nil {
		assert.False(t, bc.HashSatisfiesThreshold([]byte(dummyWordBytes), lowestThreshold))
	}

	fmt.Printf("%x\n", bc.SHA256([]byte(dummyWord)))
}
