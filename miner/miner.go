package miner

import (
	bc "blockchainnetwork/blockchain"
	"blockchainnetwork/fullnode/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Miner struct {
	// miner identification
	name string

	// the full nodes that the miner is in communication with
	nodes map[string]proto.NewFullNodeClient

	// data necessary to mine the next block
	lastBlock *bc.Block
	threshold []byte
}

func MakeMiner(name string, nodes []string) *Miner {

	miner := &Miner{
		name:  name,
		nodes: make(map[string]proto.NewFullNodeClient, len(nodes)),
	}

	// TODO: set up connections with the nodes with Dial()
	for _, node := range nodes {
		client, err := makeFullNodeClient(node)
		if err != nil {
			return nil
		}
		miner.nodes[node] = client
	}

	// TODO: start the background routine to check for new blocks in target blockchains

	return miner
}

func makeFullNodeClient(addr string) (proto.NewFullNodeClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	channel, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	return proto.NewFullNodeClient(channel), nil
}

func (miner *Miner) mine() {
	minedBlock := &bc.Block{
		Index:     miner.lastBlock.Index + 1,
		PrevHash:  miner.lastBlock.PrevHash,
		Timestamp: time.Now(),
		Nonce:     uint64(0),
		Data:      "Block mined by " + miner.name,
	}

	for !bc.HashSatisfiesThreshold(minedBlock, miner.threshold) {
		minedBlock.Nonce += 1
	}

	// send it to the blockchain

}
