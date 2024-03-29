package main

func (blockchain *Blockchain) AddBlock(data string) {
	PreviousBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	newBlock := NewBlock(data, PreviousBlock.MyBlockHash)
	blockchain.Blocks = append(blockchain.Blocks, newBlock)
}

// NewBlockChain creates a new blockchain with a genesis block
func NewBlockChain() *Blockchain {
	blockchain := &Blockchain{
		Blocks: []*Block{NewGenesisBlock()},
	}
	return blockchain
}
