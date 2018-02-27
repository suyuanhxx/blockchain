package main

import (
	//"github.com/deckarep/golang-set"
	//"container/list"
	"encoding/json"
	"crypto/sha256"
	"encoding/hex"
	"time"
	"crypto/sha1"
)

var hash = sha1.New()
var hash256 = sha256.New()

type Block struct {
	Index        int       `json:"index"`
	Timestamp    time.Time `json:"timestamp"`
	Transactions []json    `json:"transactions"`
	Proof        int       `json:"proof"`
	PreviousHash string    `json:"previous_hash"`
}

type Transaction struct {
	Amount    int    `json:"amount"`
	Recipient string `json:"recipient"`
	Sender    string `json:"sender"`
}

type Blockchain struct {
	CurrentTransactions []Transaction
	chain               []Block
	Nodes               []json
}

func (t *Blockchain) registerNode(address string) {

}

func (t *Blockchain) validChain(chain []string) {

}

func (t *Blockchain) resolveConflicts() {

}

func (t *Blockchain) newBlock(proof int, previousHash string) string {
	block := new(Block)
	block.Index = len(t.chain) + 1
	block.Timestamp = time.Now()
	block.Proof = proof
	if previousHash != "" {
		block.PreviousHash = previousHash
	} else {
		hash.Write([]byte(t.chain[len(t.chain)-1]))
		block.PreviousHash = string(hash.Sum(nil))
	}
	output, _ := json.Marshal(&block)
	return string(output)
}

func (t *Blockchain) newTransaction(sender string, recipient string, amount float64) {

}

func (t *Blockchain) lastBlock() {

}

func (t *Blockchain) hash(block []string) string {
	var blockString string
	json.Unmarshal([]byte(block), &blockString)
	hash256.Write([]byte(blockString))
	hashResult := hex.EncodeToString(hash256.Sum(nil))
	return hashResult
}

func (t *Blockchain) proof_of_work(lastProof int) int {
	proof := 0
	for !ValidProof(lastProof, proof) {
		proof += 1
	}
	return 0
}

func ValidProof(lastProof int, proof int) bool {
	return false
}

func Mine() {

}

func NewTransaction() {

}

func FullChain() {

}

func RegisterNodes() {

}

func Consensus() {

}
