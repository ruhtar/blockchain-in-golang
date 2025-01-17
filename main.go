package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var Blockchain []Block

func main() {
	godotenv.Load()

	fmt.Println("Creating Genesis Block...")
	genesisBlock := createGenesisBlock()
	fmt.Println("Generating blockchain...")
	Blockchain = append(Blockchain, genesisBlock)
	fmt.Println("Blockchain ready!")

	run()
}

type Block struct {
	Index     int    // is the position of the data record in the blockchain
	Timestamp string // is automatically determined and is the time the data is written
	Data      []byte //beats per minute
	Hash      string //is a SHA256 identifier representing this data record
	PrevHash  string //is the SHA256 identifier of the previous record in the chain
}

type Message struct {
	Data string
}

func calculateHash(block Block) string {
	data := string(block.Index) + block.Timestamp + string(block.Data) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(data))
	hash := h.Sum(nil) //nil because no other data is passed to be added to the end of the calculus
	return hex.EncodeToString(hash)
}

func createGenesisBlock() Block {
	var genesisBlock Block
	genesisBlock.Data = []byte("Genesis block")
	genesisBlock.Index = 0
	genesisBlock.Timestamp = time.Now().Format(time.RFC3339)
	genesisBlock.PrevHash = ""
	genesisBlock.Hash = calculateHash(genesisBlock)
	return genesisBlock
}

func createNewBlock(data string) (Block, error) {
	var newBlock Block
	previousBlock := Blockchain[len(Blockchain)-1]
	newBlock.Data = []byte(data)
	newBlock.Index = previousBlock.Index + 1
	newBlock.Timestamp = time.Now().Format(time.RFC3339)
	newBlock.PrevHash = previousBlock.PrevHash
	newBlock.Hash = calculateHash(newBlock)
	return newBlock, nil
}

func isBlockValid(newBlock Block, oldBlock Block) bool {
	if oldBlock.Index != newBlock.Index-1 {
		return false
	}

	if newBlock.PrevHash != oldBlock.Hash && newBlock.PrevHash != "" {
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

func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", os.Getenv("ADDR"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := createNewBlock(m.Data)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
