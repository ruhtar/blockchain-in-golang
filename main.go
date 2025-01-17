package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

var Blockchain []Block

func main() {

	fmt.Println("Creating Genesis Block...")
	genesisBlock := createGenesisBlock()
	fmt.Println("Generating blockchain...")
	Blockchain = append(Blockchain, genesisBlock)
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

func createGenesisBlock() Block {
	var genesisBlock Block
	// newBlock.BPM = bpm
	genesisBlock.Index = 0
	genesisBlock.Timestamp = time.Now().Format(time.RFC3339)
	genesisBlock.PrevHash = ""
	genesisBlock.Hash = calculateHash(genesisBlock)
	return genesisBlock
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

func isValidBlock(newBlock Block, oldBlock Block) bool {
	if oldBlock.Index != newBlock.Index+1 {
		return false
	}

	if newBlock.PrevHash != oldBlock.Hash {
		return false
	}

	if oldBlock.Timestamp >= newBlock.Timestamp {
		return false
	}

	if newBlock.Hash != calculateHash(newBlock) {
		return false
	}

	return true
}

// The blockchain uses the principle of the longest chain to validate which chain is the valid one.
func replaceChain(blocks []Block) {
	if len(blocks) > len(Blockchain) {
		Blockchain = blocks
	}
}
