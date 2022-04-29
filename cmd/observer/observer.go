package main

import (
	"context"
	"flag"
	"log"

	node "blockchainnetwork/node"
	proto "blockchainnetwork/node/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	loneNodeAddr = flag.String(
		"node-addr",
		"[::1]:8080",
		"The address for the lonely node",
	)
)

func serviceConn(address string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return grpc.Dial(address, opts...)
}

func main() {
	flag.Parse()

	nodeConn, err := serviceConn(*loneNodeAddr)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer nodeConn.Close()
	client := proto.NewNodeClient(nodeConn)

	// try getting the last block
	ctx := context.Background()
	req := &proto.GetLastBlockRequest{}

	resp, err := client.GetLastBlock(ctx, req)
	if err == nil {
		node.PrintProtoBlock(resp.LastBlock)
	}
}
