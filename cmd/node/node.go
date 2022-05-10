package main

import (
	"flag"
	"log"
	"net"

	nd "blockchainnetwork/node"
	proto "blockchainnetwork/node/proto"

	"google.golang.org/grpc"
)

var (
	nodeAddr = flag.String(
		"node-addr",
		"[::1]:8080",
		"The address for the lonely node",
	)
)

func main() {
	flag.Parse()
	nodeAddrs := flag.Args()

	log.Printf(
		"starting the node with flags: --node-addr=%s\n",
		*nodeAddr,
	)

	flag.Parse()
	lis, err := net.Listen("tcp", *nodeAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	nodePool := nd.MakeGRPCNodeClientPool(nodeAddrs)
	s := grpc.NewServer()
	proto.RegisterNodeServer(
		s,
		nd.MakeNode("tosun_node", nodePool, nodeAddrs),
	)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
