package main

import (
	"flag"
	"log"

	mnr "blockchainnetwork/miner"
	nd "blockchainnetwork/node"
)

var (
	minerName = flag.String(
		"miner-name",
		"tosun",
		"The name of the miner",
	)

	mineDelay = flag.Int(
		"mine-delay",
		0,
		"Delay coefficient for the mining process",
	)
)

func main() {
	flag.Parse()
	nodeAddrs := flag.Args()

	log.Printf(
		"starting the miner with flag: --miner-name=%s, --mine-delay=%d\n",
		*minerName,
		*mineDelay,
	)

	nodePool := nd.MakeGRPCNodeClientPool(nodeAddrs)
	miner := mnr.MakeMiner(*minerName, nodePool, nodeAddrs, uint64(*mineDelay))

	miner.MineContinuously()
}
