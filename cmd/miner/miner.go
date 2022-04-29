package main

import (
	"flag"
	"log"

	mnr "blockchainnetwork/miner"
	nd "blockchainnetwork/node"
)

var (
	loneNodeAddr = flag.String(
		"node-addr",
		"[::1]:8080",
		"The address for the lonely node",
	)
)

func main() {
	flag.Parse()

	log.Printf(
		"starting the miner with flags: --node-addr=%s\n",
		*loneNodeAddr,
	)

	// if err != nil {
	// 	log.Fatalf("fail to dial: %v", err)
	// }

	nodePool := nd.MakeGRPCNodeClientPool([]string{*loneNodeAddr})
	miner := mnr.MakeMiner("tosun_miner", nodePool)

	miner.Mine()

}