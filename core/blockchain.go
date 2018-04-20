package core

import (
	"encoding/json"
	"time"
	"fmt"
	"net/http"
	"github.com/deckarep/golang-set"
)

type Block struct {
	Index        int           `json:"index"`
	Timestamp    string        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Proof        int           `json:"proof"`
	PreviousHash string        `json:"previous_hash"`
}

type Transaction struct {
	Amount    int    `json:"amount"`
	Recipient string `json:"recipient"`
	Sender    string `json:"sender"`
}

type Chain struct {
	Chain  []Block `json:"chain"`
	Length int     `json:"length"`
}

type Blockchain struct {
	CurrentTransactions []Transaction
	Chain               []Block
	Nodes               mapset.Set
}

func (t *Blockchain) RegisterNode(host string) {
	if t.Nodes == nil {
		t.Nodes = mapset.NewSet()
	}
	t.Nodes.Add(host)
}

// Create the genesis block
func (t *Blockchain) New() *Blockchain {
	t = new(Blockchain)
	t.NewBlock(100, "1")
	return t
}

func (t *Blockchain) ValidChain(chain []Block) bool {
	lastBlock := chain[0]
	currentIndex := 1

	for ; currentIndex < len(chain); {
		block := chain[currentIndex]
		fmt.Println(lastBlock, block)
		if block.PreviousHash != Hash(lastBlock) {
			return false
		}

		if !t.ValidProof(lastBlock.Proof, block.Proof, lastBlock.PreviousHash) {
			return false
		}
		lastBlock = block
		currentIndex += 1
	}
	return true
}

func (t *Blockchain) ResolveConflicts() bool {
	var newChain []Block

	maxLength := len(t.Chain)

	for node := range t.Nodes.Iter() {
		resp, error := http.Get("http://" + node.(string) + "/chain")
		if error != nil {
			continue
		}
		defer resp.Body.Close()
		var chain Chain
		json.NewDecoder(resp.Body).Decode(&chain)

		if chain.Length > maxLength && t.ValidChain(chain.Chain) {
			maxLength = chain.Length
			newChain = chain.Chain
		}
	}

	if newChain != nil {
		t.Chain = newChain
		return true
	}
	return false

}

func (t *Blockchain) NewBlock(proof int, previousHash string) Block {
	block := new(Block)
	block.Index = len(t.Chain) + 1
	block.Timestamp = time.Now().String()
	block.Proof = proof
	if previousHash != "" {
		block.PreviousHash = previousHash
	} else {
		block.PreviousHash = Hash(t.Chain[len(t.Chain)-1])
	}

	if len(t.CurrentTransactions) != 0 {
		block.Transactions = append(block.Transactions, t.CurrentTransactions[0])
	}
	t.CurrentTransactions = []Transaction{}
	t.Chain = append(t.Chain, *block)
	return *block
}

func (t *Blockchain) NewTransaction(sender string, recipient string, amount int) int {
	transaction := new(Transaction)
	transaction.Sender = sender
	transaction.Recipient = recipient
	transaction.Amount = amount
	t.CurrentTransactions = append(t.CurrentTransactions, *transaction)
	return t.LastBlock().Index + 1
}

func (t *Blockchain) LastBlock() Block {
	return t.Chain[len(t.Chain)-1]
}