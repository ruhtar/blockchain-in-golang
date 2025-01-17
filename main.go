package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func main() {
	fmt.Println("hello world")
}

type Block struct {
	Index     int    // is the position of the data record in the blockchain
	Timestamp string // is automatically determined and is the time the data is written
	BPM       int    //beats per minute
	Hash      string //is a SHA256 identifier representing this data record
	PrevHash  string //is the SHA256 identifier of the previous record in the chain
}

func calculateHash(block Block) string {
	data := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(data))
	hash := h.Sum(nil) //nil because no other data is passed to be added to the end of the calculus
	return hex.EncodeToString(hash)
}

func createNewBlock(previousBlock Block, bpm int) (Block, error) {
	var newBlock Block
	newBlock.BPM = bpm
	newBlock.Index = previousBlock.Index + 1
	newBlock.Timestamp = time.Now().Format(time.RFC3339)
	newBlock.PrevHash = previousBlock.PrevHash
	newBlock.Hash = calculateHash(newBlock)
	return newBlock, nil
}
