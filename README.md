# BlockchainNetwork

A P2P blockchain network from scratch, written in Go.

## Installation and Usage
First, clone this repository:
```
git clone https://github.com/kilyig/BlockchainNetwork.git
```
The command line programs are in `cmd`. To start a node from the root directory, run
```
go run cmd/node/node.go --node-addr [address + port of this node] [address + port of each node that this node will communicate with]*
```
For example, if our node is at `[::1]:8080` and communicates with nodes at addresses `[::1]:8081` and `[::1]:8082`, then we would run
```
go run cmd/node/node.go --node-addr [::1]:8080 [::1]:8081 [::1]:8082
```
In this scenario, `[::1]:8080` will try to sychronize its blockchain with `[::1]:8081` and `[::1]:8082`. To start a miner from the root directory, run
```
go run cmd/miner/miner.go --miner-name [a name for this miner] --mine-delay [delay (nanoseconds) before checking a block while mining] [address + port of each node that this miner will communicate with]*
```
For example, if a miner named `tosun` stops for 100 nanoseconds before each block and communicates with nodes at addresses `[::1]:8080` and `[::1]:8081`, then we would run
```
go run cmd/miner/miner.go --miner-name tosun --mine-delay 100 [::1]:8080 [::1]:8081
```
The `--mine-delay` parameter helps us simulate miners with different hash powers.

`observer.go` can be used to check if everything is going well. It prints the last block of the node provided by the parameter `--node-addr`:
```
go run cmd/observer/observer.go --node-addr [::1]:8081
```

Note: `launch_nodes.sh` and `launch_miners.sh` might be worth checking out but they do not currently support networking.

### Using the Docker containers
In the root directory, run
```
docker build -t node -f Dockerfile.node .
docker build -t miner -f Dockerfile.miner .
```
to build the Docker containers for nodes and miners. Docker containers currently support networks with any number of miners and at most one node. To run the node, run
```
docker run --network=host --rm -p 8080:8080 node
```
To run a miner connected to this node, run
```
docker run --network=host --rm -p 8080:8080 miner --miner-name name_of_this_miner [::1]:8080
```


### Compiling the .proto code
The repo has already run the protobuf compiler to generate the proto code. If you want to do this on your own, you can do `make` to use the `Makefile` in the root directory. The `Makefile` runs the following command:
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative node/proto/node.proto
```
You will need to install `protoc`, the protobuf compiler, to get started. You can find instructions in the gPRC documentation: https://grpc.io/docs/protoc-installation/
I highly suggest you to follow the instructions under "Install pre-compiled binaries (any OS)". The `apt` and `apt-get` on Linux install an older version, which is undesirable. Additionally, you will need to install the Go plugins, following these steps: https://grpc.io/docs/languages/go/quickstart/


## Tests
The main testing method for the network was running the nodes and miners locally and checking the logs. The code comes with several basic tests for the blockchain data structure, but they do not have much coverage. To run them, navigate to `blockchain/test` and run
```
go test -run=Test -v
```

