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
)

func main() {
	flag.Parse()
	nodeAddrs := flag.Args()

	log.Printf(
		"starting the miner with flag: --miner-name=%s\n",
		*minerName,
	)

	log.Println(nodeAddrs)

	nodePool := nd.MakeGRPCNodeClientPool(nodeAddrs)
	miner := mnr.MakeMiner(*minerName, nodePool, nodeAddrs)

	miner.MineContinuously()
}
