package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

func (block *Block) SetHash() {
	block.Timestamp = time.Now().UTC().Unix()

	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	dataBytes := []byte(block.Data)
	prevHashHex := hex.EncodeToString(block.PreviousBlockHash)
	prevHashBytes := []byte(prevHashHex)

	headers := bytes.Join([][]byte{timestamp, prevHashBytes, dataBytes}, []byte{})
	hash := sha256.Sum256(headers)
	block.MyBlockHash = hash[:]

	fmt.Println("🔹 Block Created:")
	fmt.Println("   Data:", block.Data)
	fmt.Println("   Timestamp (UTC):", block.Timestamp)
	fmt.Println("   Prev Hash (Hex):", prevHashHex)
	fmt.Println("   Computed Hash (Hex):", hex.EncodeToString(block.MyBlockHash))
}

func NewBlock(data string, prevblockHash []byte) *Block {
	// Ensure Data and AllData are consistent
	dataBytes := []byte(data) // Convert to bytes once

	block := &Block{
		Data:              data,
		Timestamp:         time.Now().Unix(),
		PreviousBlockHash: prevblockHash,
		MyBlockHash:       []byte{},  // Will be set in SetHash()
		AllData:           dataBytes, // Store the exact bytes used for hashing
	}

	block.SetHash() // Compute the hash using consistent data
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
