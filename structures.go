package main

type Block struct {
	Data              string
	Timestamp         int64
	PreviousBlockHash []byte
	MyBlockHash       []byte
	AllData           []byte
}

type Blockchain struct {
	Blocks []*Block // blockchain is a series of blocks
}
