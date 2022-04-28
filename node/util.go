package node

import (
	bc "blockchainnetwork/blockchain"
	proto "blockchainnetwork/node/proto"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func ProtoBlockToBlockchainBlock(protoBlock *proto.Block) *bc.Block {
	// check if the timestamp format is convertible
	// TODO: is this needed? does it really work?
	// https://pkg.go.dev/google.golang.org/protobuf/types/known/timestamppb#hdr-Conversion_to_a_Go_Time
	if err := protoBlock.Timestamp.CheckValid(); err != nil {
		return nil
	}

	return &bc.Block{
		Index:     protoBlock.Index,
		PrevHash:  protoBlock.PrevHash,
		Timestamp: protoBlock.Timestamp.AsTime(),
		Nonce:     protoBlock.Nonce,
		Data:      protoBlock.Data,
	}
}

func BlockchainBlockToProtoBlock(blockchainBlock *bc.Block) *proto.Block {
	return &proto.Block{
		Index:     blockchainBlock.Index,
		PrevHash:  blockchainBlock.PrevHash,
		Timestamp: timestamppb.New(blockchainBlock.Timestamp),
		Nonce:     blockchainBlock.Nonce,
		Data:      blockchainBlock.Data,
	}
}
