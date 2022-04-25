.PHONY: compile

compile:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative blockchain/proto/blockchain.proto