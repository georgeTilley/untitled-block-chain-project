package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

func (block *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(block.TimeStamp, 10))
	headers := bytes.Join([][]byte{timestamp, block.PreviousBlockHash, block.AllData}, []byte{})
	hash := sha256.Sum256(headers)
	block.MyBlockHash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), prevBlockHash, []byte{}, []byte(data)}
	block.SetHash()
	return block
}

// NewGenesisBlock creates a new genesis block for the blockchain
func NewGenesisBlock() *Block {
	return &Block{
		TimeStamp:         0,
		PreviousBlockHash: []byte{},
		MyBlockHash:       []byte{},
		AllData:           []byte("Genesis Block"),
	}
}
