package core

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (blockchain *Blockchain) AddBlock(data string) error {
	db, err := Dbinit()
	if err != nil {
		return fmt.Errorf("error connecting to db: %w", err)
	}
	defer db.Close()
	PreviousBlock := blockchain.Blocks[len(blockchain.Blocks)-1]

	prevHashHex := hex.EncodeToString(PreviousBlock.MyBlockHash)

	newBlock := NewBlock(data, []byte(prevHashHex))
	blockchain.Blocks = append(blockchain.Blocks, newBlock)

	newHashHex := hex.EncodeToString(newBlock.MyBlockHash)
	newTimestamp := strconv.FormatInt(newBlock.Timestamp, 10)
	fmt.Println(" Storing Block:")
	fmt.Println("   Data:", data)
	fmt.Println("   Timestamp:", newTimestamp)
	fmt.Println("   Previous Hash (hex):", prevHashHex)
	fmt.Println("   Current Block Hash (hex):", newHashHex)

	appendfile("logs.txt", []byte(data+"\n"+newTimestamp+"\n"+prevHashHex+"\n"+newHashHex+"\n"))
	err = saveBlockToDB(db, data, newTimestamp, prevHashHex, newHashHex)
	if err != nil {
		return fmt.Errorf("error saving data to database: %w", err)
	}

	fmt.Println(" Block successfully added.")
	return nil
}

func NewBlockchain() *Blockchain {
	content, err := os.ReadFile("logs.txt")
	if err != nil {
		fmt.Println(" No existing blockchain found. Creating a new one...")
		return &Blockchain{[]*Block{NewGenesisBlock()}}
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) < 4 {
		fmt.Println(" Blockchain file incomplete. Creating a new one...")
		return &Blockchain{[]*Block{NewGenesisBlock()}}
	}

	lastData := lines[len(lines)-4]
	lastTimestamp := lines[len(lines)-3]
	lastPrevHash := lines[len(lines)-2]
	lastHash := lines[len(lines)-1]

	timestampInt, err := strconv.ParseInt(lastTimestamp, 10, 64)
	if err != nil {
		fmt.Println(" Error parsing timestamp. Starting fresh...")
		return &Blockchain{[]*Block{NewGenesisBlock()}}
	}

	prevHashBytes, _ := hex.DecodeString(lastPrevHash)
	hashBytes, _ := hex.DecodeString(lastHash)

	fmt.Println(" Loaded Last Block from File:")
	fmt.Println("   Data:", lastData)
	fmt.Println("   Timestamp:", lastTimestamp)
	fmt.Println("   Previous Hash (hex):", lastPrevHash)
	fmt.Println("   Hash (hex):", lastHash)

	lastBlock := &Block{
		Data:              lastData,
		Timestamp:         timestampInt,
		PreviousBlockHash: prevHashBytes,
		MyBlockHash:       hashBytes,
	}

	return &Blockchain{[]*Block{lastBlock}}
}
