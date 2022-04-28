package main

import (
	"flag"
	"log"

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

	log.Printf(
		"starting the miner with flags: --node=%s\n",
		*loneNodeAddr,
	)

	nodeConn, err := serviceConn(*loneNodeAddr)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer nodeConn.Close()
	//nodeClient := proto.NewNodeClient(nodeConn)

}
