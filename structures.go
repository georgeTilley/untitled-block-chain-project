package main

// Block represents a single block in the blockchain
type Block struct {
	TimeStamp         int64  `json:"time"`
	PreviousBlockHash []byte `json:"previousHash"`
	MyBlockHash       []byte `json:"myHash"`
	AllData           []byte `json:"data"`
}

// Blockchain represents a blockchain holding a list of blocks
type Blockchain struct {
	Blocks []*Block
}
