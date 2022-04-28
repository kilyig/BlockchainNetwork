package node

import (
	bc "blockchainnetwork/blockchain"
	proto "blockchainnetwork/node/proto"
)

func ProtoBlocksToBlockchainBlocks(protoBlocks []*proto.Block) []*bc.Block {
	blockchainBlocks := make([]*bc.Block, len(protoBlocks))
	for i, protoBlock := range protoBlocks {
		blockchainBlocks[i] = ProtoBlockToBlockchainBlock(protoBlock)
	}

	return blockchainBlocks
}

func BlockchainBlocksToProtoBlocks(blockchainBlocks []*bc.Block) []*proto.Block {
	protoBlocks := make([]*proto.Block, len(blockchainBlocks))
	for i, blockchainBlock := range blockchainBlocks {
		protoBlocks[i] = BlockchainBlockToProtoBlock(blockchainBlock)
	}

	return protoBlocks
}

func ProtoBlockToBlockchainBlock(protoBlock *proto.Block) *bc.Block {
	return &bc.Block{
		Index:     protoBlock.Index,
		PrevHash:  protoBlock.PrevHash,
		Timestamp: protoBlock.Timestamp,
		Nonce:     protoBlock.Nonce,
		Data:      protoBlock.Data,
	}
}

func BlockchainBlockToProtoBlock(blockchainBlock *bc.Block) *proto.Block {
	return &proto.Block{
		Index:     blockchainBlock.Index,
		PrevHash:  blockchainBlock.PrevHash,
		Timestamp: blockchainBlock.Timestamp,
		Nonce:     blockchainBlock.Nonce,
		Data:      blockchainBlock.Data,
	}
}
